import os
import sys

if len(sys.argv) != 11:
    print("~ JMP_PORT JMP_USER JMP_HOSTNAME JMP_PASSWD TARGET_PORT TARGET_USER TARGET_HOSTNAME TARGET_PASSWD PYFILE_LOCAL_PATH PYFILE_REMOTE_PATH")
    exit(-1)

JMP_PORT = sys.argv[1]
JMP_USER = sys.argv[2]
JMP_HOSTNAME = sys.argv[3]
JMP_PASSWD = sys.argv[4]
TARGET_PORT = sys.argv[5]
TARGET_USER = sys.argv[6]
TARGET_HOSTNAME = sys.argv[7]
TARGET_PASSWD = sys.argv[8]
PYFILE_LOCAL_PATH = sys.argv[9]
PYFILE_REMOTE_PATH = sys.argv[10]


CMD1 = f"sshpass -p {JMP_PASSWD} scp -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no -r -P {JMP_PORT} {PYFILE_LOCAL_PATH} {JMP_USER}@{JMP_HOSTNAME}:/tmp"
CMD2 = f"sshpass -p {JMP_PASSWD} ssh -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no -p {JMP_PORT} {JMP_USER}@{JMP_HOSTNAME} sshpass -p {TARGET_PASSWD} scp -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no -r -P {TARGET_PORT} /tmp/{os.path.basename(PYFILE_LOCAL_PATH)} {TARGET_USER}@{TARGET_HOSTNAME}:{PYFILE_REMOTE_PATH}"
CMD3 = f"sshpass -p {JMP_PASSWD} ssh -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no -p {JMP_PORT} {JMP_USER}@{JMP_HOSTNAME} sshpass -p {TARGET_PASSWD} ssh -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no -p {TARGET_PORT} {TARGET_USER}@{TARGET_HOSTNAME} python3 {PYFILE_REMOTE_PATH}"

os.system(CMD1)
os.system(CMD2)
os.system(CMD3)