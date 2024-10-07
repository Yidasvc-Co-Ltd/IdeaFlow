#!/usr/bin/python3
import os
import sys
import signal
import socket
import sqlite3
import getopt
import subprocess
    
try:
    from termcolor import colored
    from flask import Flask, request, jsonify
    import datetime
    import requests
    import threading
except ImportError:
    os.system('''
            pip3 install termcolor
            pip3 install flask
            pip3 install datetime
            pip3 install requests
            pip3 install threading
            ''')
    from termcolor import colored
    from flask import Flask, request, jsonify
    import datetime
    import requests
    import threading


LOG_ERROR = 0
LOG_WARNING = 1
LOG_INFO = 2
LOG_VERBOSE = 3

APPLICATION_TYPE = "Dynamic Task Daemon Service"
DB_FILE = None
IF_VERBOSE = None
SERVICE_PORT = None

conn = None
c = None
flask_app = None
service_request_counter = 0
running_task={
    "documentID":None,
    "userID":None,
    "dynFigureID":None,
    "taskName":None
}

def log(msg, component, LEVEL):
    def log_uncolored_str(msg, component, level_string):
        return f"[{datetime.datetime.now().strftime('%Y-%m-%d %H:%M:%S')}][{level_string:>7}][{component:>15}] {msg}"
    
    if LEVEL==LOG_ERROR:
        print(colored(log_uncolored_str(msg, component, "ERROR"),"red"))
        exit_safe(-1)
    elif LEVEL==LOG_WARNING:
        print(colored(log_uncolored_str(msg, component, "WARNING"),"yellow"))
    elif LEVEL==LOG_INFO:
        print(colored(log_uncolored_str(msg, component, "INFO"),"green"))
    elif LEVEL==LOG_VERBOSE:
        if IF_VERBOSE==1:
            print(colored(log_uncolored_str(msg, component, "VERBOSE"), "cyan"))
            

def read_args():
    global DB_FILE
    global IF_VERBOSE
    global SERVICE_PORT
    
    opts, args = getopt.getopt(sys.argv[1:], "", ["dbfile=","verbose","port="])
    for opt_name,opt_value in opts:
        if opt_name=="--dbfile":
            DB_FILE = opt_value
            log(f"DB_FILE = {DB_FILE} (getopt)", "READ_ARGS", LOG_INFO)
        elif opt_name=="--verbose":
            IF_VERBOSE = 1
            log(f"IF_VERBOSE = {IF_VERBOSE} (getopt)", "READ_ARGS", LOG_INFO)
        elif opt_name=="--port":
            SERVICE_PORT = int(opt_value)
            log(f"SERVICE_PORT = {SERVICE_PORT} (getopt)", "READ_ARGS", LOG_INFO)
        
    if IF_VERBOSE == None:
        IF_VERBOSE = 0
        log(f"IF_VERBOSE = {IF_VERBOSE} (default)", "READ_ARGS", LOG_INFO)
    if DB_FILE == None:
        home_dir = os.popen('echo $HOME').readlines()[0]
        if len(home_dir)>1:
            home_dir = home_dir[:-1]
        DB_FILE = os.path.join(home_dir,".dynTaskDB.db")
        log(f"DB_FILE = {DB_FILE} (default)", "READ_ARGS", LOG_INFO)
    if SERVICE_PORT == None:
        SERVICE_PORT = 1231
        log(f"SERVICE_PORT = {SERVICE_PORT} (default)", "READ_ARGS", LOG_INFO)
        
        
def detect_port():
    global SERVICE_PORT
    global APPLICATION_TYPE
    def check_port_exist():
        sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)  
        sock.settimeout(5)  
        try:  
            sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
            sock.settimeout(1)
            sock.connect(("127.0.0.1", SERVICE_PORT))
            sock.shutdown(socket.SHUT_RDWR)
            log(f"Socket exist on port {SERVICE_PORT}", "DETECT_PORT", LOG_VERBOSE)
            return True
        except Exception as e:
            log(str(e), "DETECT_PORT", LOG_VERBOSE)
            log(f"Socket not exist on {SERVICE_PORT}", "DETECT_PORT", LOG_VERBOSE)
            return False
            
    
    if check_port_exist()==True:
        try:
            res = requests.post(f"http://127.0.0.1:{SERVICE_PORT}/", json={"op":"heartbeat"}).json()
            if res["state"] == 200 and res["application"] == APPLICATION_TYPE:
                log(f"TARGET=AVAILABLE", "DETECT_PORT", LOG_INFO)
                exit_safe()
        except Exception as e:
            log(str(e), "DETECT_PORT", LOG_WARNING)
            log(f"Binding service port failed. Please check if port {SERVICE_PORT} is bound by other application. Exit.", "DETECT_PORT", LOG_ERROR)
        
    log("No application instance running and service port is available. Initing application ...", "DETECT_PORT", LOG_VERBOSE)

            
def sqlite_init():
    global conn
    global c

    conn = sqlite3.connect(DB_FILE, check_same_thread=False)
    c = conn.cursor()
    try:
        log("Trying to select * from Tasks ...", "SQLITE_INIT", LOG_VERBOSE)
        c.execute('''
            select * from Tasks
            ''')
    except sqlite3.OperationalError:
        log("Failed to select * from Tasks. Creating table Tasks ...", "SQLITE_INIT", LOG_VERBOSE)
        c.execute('''create table Tasks(
            documentID varchar not null,
            userID varchar not null,
            dynFigureID varchar not null,
            taskName varchar not null,
            codeShell longtext,
            timeStart datetime,
            timeEnd datetime,
            tag varchar,
            stdout longtext,
            stderr longtext,
            returnCode integer,
            primary key (documentID, userID, dynFigureID, taskName)
            );''')
        c.execute('''create table RunLocks(
            lockHandler integer primary key,
            lockState integer
            );''')
    conn.commit()
    RunLock_off()
    c.execute("update Tasks set timeStart=null where timeEnd is null")
    conn.commit()
    threading.Thread(target=task_add_hook).start()

    
def RunLock_on():
    global conn
    global c
    results = c.execute("select * from RunLocks where lockHandler=0").fetchall()
    if len(results)==0:
        c.execute("insert into RunLocks values(0,0)")
        conn.commit()
    c.execute("update RunLocks set lockState=1 where lockHandler=0")
    conn.commit()
    
def RunLock_test():
    global conn
    global c
    results = c.execute("select * from RunLocks where lockHandler=0").fetchall()
    if len(results)==0:
        c.execute("insert into RunLocks values(0,0)")
        conn.commit()
        return False
    results = c.execute("select lockHandler,lockState from RunLocks where lockHandler=0").fetchall()
    if int(results[0][1])==1:
        return True
    else:
        return False
    
def RunLock_off():
    global conn
    global c
    results = c.execute("select * from RunLocks where lockHandler=0").fetchall()
    if len(results)==0:
        c.execute("insert into RunLocks values(0,0)")
        conn.commit()
    c.execute("update RunLocks set lockState=0 where lockHandler=0");
    conn.commit()

def task_add_hook():
    global conn
    global c    
    global running_task
    results = c.execute("select documentID, userID, dynFigureID, taskName, codeShell from Tasks where timeStart is null limit 1").fetchall()
    if len(results)>0:
        documentID = results[0][0]
        userID = results[0][1]
        dynFigureID = results[0][2]
        taskName = results[0][3]
        codeShell = results[0][4]
        if RunLock_test()==True:
            return
        
        RunLock_on()
        running_task["documentID"] = documentID
        running_task["userID"] = userID
        running_task["dynFigureID"] = dynFigureID
        running_task["taskName"] = taskName
        timeStart = datetime.datetime.now().strftime(r"%Y-%m-%d%H:%M:%S")
        c.execute("update Tasks set timeStart=? where documentID=? and userID=? and dynFigureID=? and taskName=?", (timeStart, documentID, userID, dynFigureID, taskName))
        conn.commit()
        log(f"Start task : {codeShell}","EXECUTOR",LOG_INFO)
        try:
            p = subprocess.Popen(codeShell.split(), stdin=subprocess.PIPE, stdout=subprocess.PIPE, stderr=subprocess.PIPE, shell=False)
            log(f"Finish task : {codeShell}","EXECUTOR",LOG_INFO)
            stdout, stderr = p.communicate()
            stdout = str(stdout, encoding="utf-8")
            stderr = str(stderr, encoding="utf-8")
            returnCode = p.poll()
            timeEnd = datetime.datetime.now().strftime(r"%Y-%m-%d%H:%M:%S")
            taskName = running_task["taskName"]
            c.execute("update Tasks set timeEnd=?,stdout=?,stderr=?,returnCode=? where documentID=? and userID=? and dynFigureID=? and taskName=?", (timeEnd, stdout, stderr, returnCode, documentID, userID, dynFigureID, taskName))
            conn.commit()
        except Exception as e:
            log(f"Error task : {codeShell} : {str(e)}", "EXECUTOR", LOG_WARNING)
            timeEnd = datetime.datetime.now().strftime(r"%Y-%m-%d%H:%M:%S")
            taskName = running_task["taskName"]
            c.execute("update Tasks set timeEnd=? where documentID=? and userID=? and dynFigureID=? and taskName=?", (timeEnd, documentID, userID, dynFigureID, taskName))
            conn.commit()
        RunLock_off()
        
        task_add_hook()
        
        
    
    
def task_rename_hook(documentID, userID, dynFigureID, oldName, newName):
    global running_task
    if documentID==running_task.documentID and userID==running_task.userID and dynFigureID==running_task.dynFigureID and oldName==running_task.taskName:
        running_task.taskName = newName
        

def service_init():
    global APPLICATION_TYPE
    global flask_app
    
    flask_app = Flask(APPLICATION_TYPE)
    debug = True if IF_VERBOSE else False
    
    @flask_app.route("/", methods=["POST"])
    def flask_service_handler():
        global service_request_counter
        req = request.json
        res = {}
        log(f"req[#{service_request_counter}] = {req}", "FLASK", LOG_VERBOSE)
        op = req.get("op")
        
        try:
            if op == "heartbeat":
                res["state"] = 200
                res["application"] = APPLICATION_TYPE
            elif op == "query":
                documentID = req.get("documentID")
                userID = req.get("userID")
                dynFigureID = req.get("dynFigureID")
                c.execute("select * from Tasks where documentID=? and userID=? and dynFigureID=?", (documentID, userID, dynFigureID))
                # print(c.fetchall())
                res["state"] = 200
                res["results"] = c.fetchall()
            elif op == "add":
                tasks = req.get("tasks")
                for i in range(len(tasks)):
                    documentID = tasks[i].get("documentID")
                    userID = tasks[i].get("userID")
                    dynFigureID = tasks[i].get("dynFigureID")
                    taskName = tasks[i].get("taskName")
                    codeShell = tasks[i].get("codeShell")
                    tag = tasks[i].get("tag")
                    c.execute("insert into Tasks (documentID, userID, dynFigureID, taskName, codeShell, tag) values(?,?,?,?,?,?)", (documentID, userID, dynFigureID, taskName, codeShell, tag))
                conn.commit()
                res["state"] = 200
                task_add_hook()
            elif op == "delete":
                keys = req.get("keys")
                for i in range(len(keys)):
                    documentID = keys[i].get("documentID")
                    userID = keys[i].get("userID")
                    dynFigureID = keys[i].get("dynFigureID")
                    taskName = keys[i].get("taskName")
                    c.execute("delete from Tasks where documentID=? and userID=? and dynFigureID=? and taskName=?", (documentID, userID, dynFigureID, taskName))
                conn.commit()
                res["state"] = 200
            elif op == "rename":
                documentID = req.get("documentID")
                userID = req.get("userID")
                dynFigureID = req.get("dynFigureID")
                oldName = req.get("oldName")
                newName = req.get("newName")
                c.execute("update Tasks set taskName=? where documentID=? and userID=? and dynFigureID=? and taskName=?", (newName, documentID, userID, dynFigureID, oldName))
                conn.commit()
                res["state"] = 200
                task_rename_hook(documentID, userID, dynFigureID, oldName, newName)
                
        except Exception as e:
            res["state"] = 500
            res["error"] = str(e)
            log(f"req[#{service_request_counter}] = {req}", "FLASK", LOG_WARNING)
            log(str(e), "FLASK", LOG_WARNING)
                
        log(f"res[#{service_request_counter}] = {res}", "FLASK", LOG_VERBOSE)
        service_request_counter += 1
        return res
    
    log(f"TARGET=AVAILABLE", "DETECT_PORT", LOG_INFO)
    flask_app.run(port=SERVICE_PORT, host="127.0.0.1", debug=debug, use_reloader=False, threaded=True)
        
    

def exit_safe(exit_value:int=0, handle_signal=None):
    if handle_signal:
        log(f"Handling signal {handle_signal}", "BEFORE_EXIT", LOG_INFO)

    try:
        if conn:
            RunLock_off()
            conn.close()
    except:
        log("Failed to free runlock or close connection.", "BEFORE_EXIT", LOG_WARNING)
        exit_value = -1
    
    if exit_value:
        exit(exit_value)
    else:
        exit(0)

def main():
    signal.signal(signal.SIGINT, lambda signum, frame:exit_safe(0, signum))

    log("read_args", "MAIN", LOG_INFO)
    read_args()
    
    log("detect_port", "MAIN", LOG_VERBOSE)
    detect_port()
    
    log("sqlite_init", "MAIN", LOG_VERBOSE)
    sqlite_init()
    
    log("service_init", "MAIN", LOG_VERBOSE)
    service_init() # 阻塞
    
    
    
    log("before_exit", "MAIN", LOG_VERBOSE)
    exit_safe()
    

if __name__ =="__main__":
    main()