
cset("import_list", []);
cset("import_reference_map", {"js/utils/import/import.js":{"js/utils/cookie/cookie.js":1}});

// function register_js(src) {
//     // //从cookie里获取key=import_list的数组。如果src存在，直接返回。否则将参数src加入import_list。
//     // let list = cget('import_list')
//     // for (i = 0; i < list.length; i++) {
//     //     if (list[i] == src) {
//     //         return 0
//     //     }
//     // }
//     // list.push(src)
//     // console.log("register_js " + src)
//     // return 1

//     // //数组有值返回0，没值返回1
//     // //return
// }

// function import_js(src) {
//     //document.write并将参数src加入import_list。
//     let res = register_js(src)
//     //res为1，执行引入功能；res为0，什么也不执行
//     if (res == 1) {
//         document.write("<script src=" + src + "></script>");
//         console.log("import_js " + src)
//     }
// }

function walk_reference_tree(map){
    // let str = ``;
    let list = [];
    let bfs_queue = [];
    for(let x in map){
        list.push(x);
        for(let y in map[x]){
            if(y in list){
                continue;
            }
            bfs_queue.push(y);
        }
    }
    while(bfs_queue.length){
        let current = bfs_queue[0];
        if(current in list){
            continue;
        }
        bfs_queue = bfs_queue.slice(1);
        list.push(current);
        for(let y in map[current]){
            if(y in list){
                continue;
            }
            bfs_queue.push(y);
        }
    }
    cset("import_list",list);
    for(let i=list.length-1;i>=0;i--){
        // str+=`<script src="${list[i]}"></script>`
        document.createElement("script");
        document.body.appendChild(`<script src="${list[i]}"></script>`);
    }
    // return str;
}

function import_js(want, current){
    let list = cget('import_list');
    if(want in list){
        return;
    }
    let map = cget("import_reference_map");
    if(!map[current]){
        map[current] = {};
    }
    map[current][want] = 1;
    console.log(map);
    cset("import_reference_map", map);
    let update_dom_worker = ()=>{
        if(window.import_mutex_lock == 0 || window.import_mutex_lock == undefined){
            window.import_mutex_lock = 1;
            document.body.innerHTML = document.querySelector("#main").outerHTML + walk_reference_tree(map);
            window.import_mutex_lock = 0;
        }else{
            setTimeout(update_dom_worker,1000);
        }
    }
    update_dom_worker();
}
