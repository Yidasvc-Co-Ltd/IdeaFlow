function render_home() {
    let data = cget("home");
    let documents_searched = (() => {
        if (data.searchbox_text.length > 0) {
            return data.documents.filter(item => item.path.toLowerCase().indexOf(data.searchbox_text.toLowerCase()) >= 0);
        } else {
            return data.documents;
        }
    })();
    let collectedDocuments = (() => {
        // return data.documents.filter(document => document.isCollected);
        console.log(data.documents);
    })();
    let textColor = (() => {
        return data.searchbox_text.length > 0 ? 'path-red' : 'path-black';
    })();


    if (data.activeButton == "all") {
        document.querySelector("#button_all").className = "button active";
        document.querySelector("#button_collection").className = "button";
    } else if (data.activeButton == "collection") {
        document.querySelector("#button_all").className = "button";
        document.querySelector("#button_collection").className = "button active";
    }

    let global_context = "";
    if (data.show) {
        // 这里是全部文件
        if (documents_searched) {
            for (let i = 0; i < documents_searched.length; i++) {
                let document_item = documents_searched[i];
                global_context += `
                <div class="context-document-item">
                    <div class="ideaFlow" onclick='javascript:push_edit_page(${JSON.stringify(document_item)})'></div>
                    <div path="${document_item.path}" class="document_path ${textColor}" >${document_item.path}</div>
                    <div class="update-time">${document_item.updateTime}</div>
                    ${document_item.showDeleteButton ?
                        `<button class="delete-button" onclick='javascript:deleteDocument(event, ${JSON.stringify(document_item)})'>X</button>` : ""}
                    ${document_item.isRenaming ?
                        `<button class="rename-button" onclick='javascript:startRenaming(event, ${JSON.stringify(document_item)})'>R</button>` : ""}
                    <button
                        class="button_star collection-button margin-right-130px ${document_item.isCollected ? 'collected' : ''}"
                        onclick='javascript:toggleCollection(event,${JSON.stringify(document_item)});'> 
                        ${document_item.isCollected ? '取消收藏' : '收藏'}
                    </button>
                </div>`
            }
        }

    } else {
        // 这是收藏文件
        for (let i = 0; i < collectedDocuments.length; i++) {
            let document_item = collectedDocuments[i];
            global_context += `
            <div class="context-document-item">
                <div class="ideaFlow" onclick='javascript:push_edit_page(${JSON.stringify(document_item)})'></div>
                <div path="${document_item.path}" class="document_path ${textColor}">${document_item.path}</div>
                <div class="update-time">${document_item.updateTime}</div>
                <button onclick='javascript:toggleCollection(event, ${JSON.stringify(document_item)});'
                    class="button_star collection-button ${document_item.isCollected ? 'collected' : ''}">
                    ${document_item.isCollected ? '取消收藏' : '收藏'}
                </button>
                ${document_item.showDeleteButton ?
                    `<button class="delete-button" onclick='javascript:deleteDocument(event, ${JSON.stringify(document_item)})'>X</button>` : ""}
                ${document_item.isRenaming ?
                    `<button class="rename-button" onclick='javascript:startRenaming(event, ${JSON.stringify(document_item)})'>R</button>` : ""}
            </div>`
        }
    }
    // 右键菜单
    global_context += `
    <transition name="fade" appear>
        ${data.isContextMenuVisible ? `
            <div class="context-menu"
                style="{ top: ${data.contextMenuPosition.y}px, left: ${data.contextMenuPosition.x}px }">
                <ul>
                    <li onclick="javascript:createNewDocument()">新建</li>
                    <li onclick="javascript:startRenaming(event, ${data.selectedDocument})">重命名</li>
                    <li onclick="javascript:deleteDocument(event, ${data.selectedDocument})">删除</li>
                </ul>
            </div>`: ""}
    </transition>`
    document.querySelector("#global-context").innerHTML = global_context;

}


function showAll() {
    let _ = cget("home");
    _.show = true;
    _.activeButton = 'all';
    cset("home", _);
    render_home();
}

function showCollection() {
    let _ = cget("home");
    _.show = false;
    _.activeButton = 'collection';
    cset("home", _);
    render_home();
}

function push_edit_page(document_info) {
    cset("fastpaper-document-info", document_info);
    window.location.href = "/edit.html";
    const data = {
        "userID": "qqq",
        "documentID": document_info.documentID,
        "operate": "Document_update_time",
        "operate_type": "Operate_documents"
    };
    axios.post('/api/backend', data)
        .then(response => {
            console.log(response);
        })
        .catch(error => {
            console.error(error);
        });
    render_home();
}

//新建path
async function createNewDocument() {
    const newpath = prompt("请输入新文档的标题", "New Document");
    // 如果用户点击了取消按钮或者没有输入标题，则不进行任何操作
    if (newpath === null || newpath === "") {
        return;
    }

    let _ = cget("home");
    // 检查是否存在重复的名称
    // const isDuplicate = _.documents.some(doc => doc.path === newpath);
    // if (isDuplicate) {
    //     alert("文档标题重复，请重新输入一个名称");
    //     return;
    // }

    // 这里确保当新建时，不会出现红色×和绿色R

    // _.documents.forEach(document => {
    //     document.showDeleteButton = false;
    //     document.isRenaming = false;
    // });
    if (_.documents) {
        for (let i = 0; i < _.documents.length; i++) {
            if (newpath == _.documents[i]) {
                alert("文档标题重复，请重新输入一个名称");
                return;
            }
            else {
                _.documents[i].showDeleteButton = false;
                _.documents[i].isRenaming = false;
            }
        }
    }

    cset("home", _);
    render_home();

    const data = {
        "userID": "qqq",
        "operate": "Document_create",
        "path": newpath,
        "operate_type": "Operate_documents"
    };
    try {
        const response = await axios.post('/api/backend', data);
        console.log('新建成功');
        await updateDocuments();
    } catch (error) {
        console.error(error);
    }
}

// 出现红色X
function toggleDeleteButton() {
    let _ = cget("home");
    _.documents.forEach(document => {
        document.showDeleteButton = !document.showDeleteButton;
    });

    //确保删除时不会出现绿色R（重命名）
    _.documents.forEach(document => {
        document.isRenaming = false;
    });
    cset("home", _);
    render_home();
}

//删除文档
async function deleteDocument(event, document) {
    event.stopPropagation(); // 阻止事件冒泡
    let _ = cget("home");
    _.isContextMenuVisible = false;
    // 显示确认对话框
    const isConfirmed = window.confirm('您确定要删除此文档吗？');

    if (isConfirmed) {
        const index = _.documents.indexOf(document);
        if (index > -1) {
            _.documents.splice(index, 1);
        }
        cset("home", _);
        render_home();

        const data = {
            "userID": "qqq",
            "operate": "Document_delete",
            "documentID": document.documentID,
            "operate_type": "Operate_documents"
        };

        try {
            await axios.post('/api/backend', data);
            console.log('删除成功');
            await updateDocuments();
        } catch (error) {
            console.error(error);
        }
    }
    render_home();
}

//出现绿色R
function toggleRenameButton() {
    let _ = cget("home");
    _.documents.forEach(document => {
        document.isRenaming = !document.isRenaming;
    });
    //这里确保重命名不会出现红色×（删除）
    _.documents.forEach(document => {
        document.showDeleteButton = false;
    });
    cset("home", _);
    render_home();
}

//重命名文档
async function startRenaming(event, document) {
    event.stopPropagation();
    let _ = cget("home");
    _.isContextMenuVisible = false;
    cset("home", _);
    render_home();
    const newpath = prompt("请输入新文档的标题", document.path);
    if (newpath === null || newpath === "") {
        return;
    }
    const data = {
        "userID": "qqq",
        "documentID": document.documentID,
        "path": newpath,
        "operate": "Document_update",
        "operate_type": "Operate_documents"
    };

    try {
        const response = await axios.post('/api/backend', data);
        console.log('重命名成功');
        await updateDocuments(); // 更新文档列表
    } catch (error) {
        console.error(error);
    }
    render_home();
}


//收藏按钮
function toggleCollection(event, document) {
    event.stopPropagation();
    document.isCollected = !document.isCollected;
    const data = {
        "userID": "qqq",
        "documentID": document.documentID,
        "isCollected": document.isCollected ? 1 : 0,
        "operate": "Document_update_is_collected",
        "operate_type": "Operate_documents"
    };
    axios.post('/api/backend', data)
        .then(response => {
            console.log('收藏成功');
            axios.post('/api/backend', {
                "userID": "qqq",
                "operate": "Document_query_all",
                "operate_type": "Operate_documents"
            }).then(response => {
                console.log(response.data);
                let _ = cget("home");
                _.documents = response.data.data;
                cset("home", _);
                render_home();
            }).catch(error => {
                console.error(error);
            });
        })
        .catch(error => {
            console.error(error);
        });

}

//搜索框
function fun(e) {
    let _ = cget("home");
    _.searchbox_text = e.target.value;
    cset("home", _);
    render_home();
}
function updateDocuments() {
    const data = {
        "userID": "qqq",
        "operate": "Document_query_all",
        "operate_type": "Operate_documents"
    };
    axios.post('/api/backend', data)
        .then(response => {
            console.log(response.data);
            let _ = cget("home");
            _.documents = response.data.data;
            cset("home", _);
            render_home();
        })
        .catch(error => {
            console.error(error);
        });
}
function showContextMenu(event, doc) {
    let _ = cget("home");
    _.selectedDocument = doc;
    _.isContextMenuVisible = true;
    _.contextMenuPosition = { x: event.clientX, y: event.clientY };
    cset("home", _);
    render_home();
    document.addEventListener('click', hideContextMenu);
    window.addEventListener('click', () => {
        let _ = cget("home");
        _.contextMenu.visible = false;
        cset("home", _);
        render_home();
    }, { once: true });
    render_home();
}
function hideContextMenu() {
    let _ = cget("home");
    _.isContextMenuVisible = false;
    cset("home", _);
    document.removeEventListener('click', hideContextMenu);
    render_home();
}



function home_mounted() {
    cset("home", {
        show: true,
        activeButton: 'all',
        documents: [],
        searchbox_text: "",
        isContextMenuVisible: false,
        contextMenuPosition: { x: 0, y: 0 },
        selectedDocument: null
    });

    axios.post('/api/backend', {
        "userID": "qqq",
        "operate": "Document_query_all",
        "operate_type": "Operate_documents"
    })
        .then(response => {
            console.log(response.data.data);
            let _ = cget("home");
            _.documents = response.data.data;
            cset("home", _);
            render_home();
        })
        .catch(error => {
            console.error(error);
        });
    render_home();
}
