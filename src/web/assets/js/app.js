//the frontend chunksize must equal backend chunksize
var ChunkSize = 4*1024*1024

Vue.prototype.$axios = axios
var vm =new Vue({
    delimiters: ['$', '$'],//use the '$$' delimiters,because the '{{' confilct with 'iris' mvc
    el: '#app',
    data:{
        tableheight:"auto",
        shareTableHeight:"auto",
        tableData: [],
        multipleSelection: [],
        uploadTableData:[],
        drawer:false,
        curRightRow:{},
        parentDir:0,
        dialogFormVisible:false,
        dir_name:"",
        DeletedialogVisible:false,
        NavArray:[{
            "ID":0,
            "FileName":"首页"
        }],
        newFileName:"",
        newFileNameDialogFormVisible:false,
        moveFileDialogFormVisible:false,
        moveFileData:[],
        moveFileTreeDefaultExpanded:[],
        moveFiledefaultProps: {
            children: 'children',
            label: 'label'
        },
        curMoveFileTreeSelected:{},
        shareDialogFormVisible:false,
        shareForm: {
            name: '分享文件(夹):',
            expdate: '7',
            haspwd: true,
            pwd: '',
            link: '',
            share_copy_text:''
        },
        share:{
            tableData:[]
        },
        createShareButtonText:"创建链接",
        mainMenuIndex:1,
        shareTableData:[],
        appHost:""

    },
    methods: {
        mouseleftclick(){
            var menu = document.querySelector("#context-menu");
            menu.style.display = 'none';
            this.curRightRow = {}
        },
        getHeight(){
            this.tableheight=window.innerHeight-121.511+'px';  //获取浏览器高度减去顶部导航栏
            this.shareTableHeight=window.innerHeight-60.511+'px';  //获取浏览器高度减去顶部导航栏


        },
        handleSelectionChange(val) {
            this.multipleSelection = val;
        },
        UserCommandHandler(command){
            switch (command) {
                case 'exit':
                    this.logout()
                    break
                case 'info':
                    vm.$message({
                        message: '这个功能还没做。。。。。。',
                        type: 'success'
                    });
                    break
            }

        },
        logout(){

            this.$axios.get('/login/logout').then(function (response) {
                if (response.data.Msg == 'OK'){

                    window.location.href= "/login"
                }
            }).catch(function (err) {
                console.log(err)
            })
        },
        table_expdate_formatter(row, column){
            if (row.ShareTime==1){
                return "1 天"

            }else if(row.ShareTime==7){
                return "7 天"

            }else if(row.ShareTime==0){
                return "永久"
            }
        },
        table_silze_formatter(row, column){
            if (parseInt(row.FileSize) == 0){
                return ''
            }else {
                if (parseInt(row.FileSize)>=1024*1024*1024){
                    Math.round( parseFloat(row.FileSize)/(1024*1024*1024)* 10) / 10
                    return  Math.round( parseFloat(row.FileSize)/(1024*1024*1024) * 10) / 10+'G'
                }else if (parseInt(row.FileSize) >= 1024*1024){
                    return Math.round( parseFloat(row.FileSize)/(1024*1024) * 10) / 10+'M'
                }else if (parseInt(row.FileSize) >= 1024){
                    return Math.round( parseFloat(row.FileSize)/(1024) * 10) / 10+'K'
                }else {
                    return row.FileSize+'b'
                }
            }
        },
        table_date_formatter(row, column){
            _obj  = row[column.property].toString().split(':')
            return  row[column.property].toString().replace(":"+_obj[_obj.length-1],"")
        },
        fileOnclick(data){
           if (data.row.IsDir){
                //go in dir
               this.parentDir = data.row.ID
               vm.NavArray.push(data.row)
               GetFiles(data.row.ID)
           }else {
               //maybe pre-show file?
           }
        },
        preDir(data){
                //go in dir
            if (data.ID ==0){
                this.NavArray =[{
                    "ID":0,
                    "FileName":"首页"
                }]
                this.parentDir = 0
            }else if (this.NavArray.length==1){
                this.parentDir = 0
                return;
            }else {
                index =vm.NavArray.indexOf(data)
                if (index<0){
                    return
                }
                this.parentDir = vm.NavArray[index].ID

                vm.NavArray.splice(index+1,(vm.NavArray.length)-index)

            }
            GetFiles(this.parentDir)




        },
        menuSelected(index,indexPath){
            this.mainMenuIndex = index
            if(index==1){
                this.NavArray =[{
                    "ID":0,
                    "FileName":"首页"
                }]
                this.parentDir = 0
                GetFiles(0)
            }else if(index==2){
                GetShareFiles()
            }else if(index==3){

            }


        },
        fileDeleteConfirm(){
            RightMenuDisplayNone()
            this.DeletedialogVisible=true
        },
        fileDelete(){
            RightMenuDisplayNone()
            this.DeletedialogVisible=false
            this.$axios.get("/file/delete/"+this.curRightRow.FileQetag+"/"+this.curRightRow.ID
            ).then(resp=>{
                if (resp.data.Status == 1){
                    GetFiles(this.curRightRow.ParentDir)
                }else {
                   ErrMsg(resp.data.Msg)
                }
            })
        },
        uploadDelete(data){
            data.row.cancel()
            index =this.uploadTableData.indexOf(data.row)
            this.uploadTableData.splice(index,1)

        },
        uploadPause(data){
            if (data.row.isUploading()){
                data.row.pause()
            }else {
                data.row.retry()
            }
        },
        row_contextmenu(row, column, event) {
            this.curRightRow = row
            var menu = document.querySelector("#context-menu");
            event.preventDefault();
            if (event.clientY+171 > window.innerHeight){
                menu.style.top = event.clientY -171 + 'px';
            }else {
                menu.style.top = event.clientY  + 'px';
            }
            menu.style.left = event.clientX + 'px';


            menu.style.display = 'block';
        },
        fileDownload(){
            RightMenuDisplayNone()
            if (this.curRightRow.IsDir==1){
                ErrMsg("暂不支持文件夹下载功能")
                return
            }
            try {
                var elemIF = document.createElement("iframe");
                elemIF.src = "/file/downloadfile/"+this.curRightRow.FileName+"?fileqetag="+this.curRightRow.FileQetag;
                elemIF.style.display = "none";
                document.body.appendChild(elemIF);
            } catch (e) {
                ErrMsg(e)
            }
        },
        OnCreateDir() {
            this.dialogFormVisible = false
           // alert(this.dir_name)

            this.$axios.get("/file/createdir/"+vm.$data.parentDir+"/"+vm.$data.dir_name).
            then(resp=>{
                if (resp.data.Status ==1){
                    vm.$data.dir_name = ""
                    GetFiles(vm.$data.parentDir)
                }else {
                    ErrMsg(resp.data.Msg)
                }
            })
        },
        preOnRenameFile(){
            this.newFileNameDialogFormVisible = true
            RightMenuDisplayNone()
        },
        OnRenameFile() {
            this.newFileNameDialogFormVisible = false
            this.$axios.get("/file/renamefile/"+vm.$data.curRightRow.ID+"/"+vm.$data.newFileName).
            then(resp=>{
                if (resp.data.Status ==1){
                    vm.$data.newFileName = ""
                    GetFiles(vm.$data.parentDir)
                }else {
                    ErrMsg(resp.data.Msg)
                }
            })
        },
        preOnMoveFile(){
            RightMenuDisplayNone()
            this.$axios.get("/file/userdirs/"+this.curRightRow.ID).
            then(resp=>{
                if (resp.data.Status ==1){
                    console.log(resp.data.Data)
                    this.moveFileDialogFormVisible =true



                    tmp=[{"id":0,"label":"全部文件","children":resp.data.Data}]
                    this.moveFileData=tmp
                    this.moveFileTreeDefaultExpanded.push(0)


                }else {
                    ErrMsg(resp.data.Msg)
                }
            })
        },
        moveFileTreeSelect(data,node,obj){
            this.curMoveFileTreeSelected= data.id
        },
        OnMoveFile(){
            this.moveFileDialogFormVisible=false
            data =new FormData()
            data.append("id",this.curRightRow.ID)
            data.append("dir",this.curMoveFileTreeSelected)

            this.$axios.post("/file/movefile",data).
            then(resp=>{
                if (resp.data.Status ==1){
                    GetFiles(this.parentDir)
                }else {
                    ErrMsg(resp.data.Msg)
                }
            })
        },
        moveFileFilterNode(value, data){
            if (data.id==value){
                return false
            }
            return true
        },
        preShareFile(){
            if(this.curRightRow.IsDir==1){
                ErrMsg("抱歉，暂时不支持文件夹分享")
                return
            }
            this.createShareButtonText="创建链接"
            RightMenuDisplayNone()
            this.shareForm.pwd = randomCode()
            this.shareForm.link = ""
            this.shareDialogFormVisible = true
            this.shareForm.name='分享文件(夹): '+this.curRightRow.FileName
        },
        OnShareFile(){

            if(this.shareForm.link.length>0){




            }else {
                data = new FormData()
                data.append("user_file_id",this.curRightRow.ID)
                data.append("share_pwd",this.shareForm.pwd)
                data.append("share_time",this.shareForm.expdate)


                this.$axios.post("/share/createshare",data)
                    .then(resp=>{
                        if (resp.data.Status ==1){
                            this.shareForm.link =window.location.href+"share/"+resp.data.Data.link
                            this.shareForm.share_copy_text="打开链接:"+this.shareForm.link+"  密码:"+this.shareForm.pwd +" 查看我分享给你的文件。"
                            this.shareDialogFormVisible = true
                        }else {
                            ErrMsg(resp.data.Msg)
                        }
                    }).catch(err=>{
                    ErrMsg(err)
                })
            }


        },
        cancelShare(props){
            this.$axios.get("/share/cancelshare/"+props.row.ShareId)
                .then(resp=>{
                    if (resp.data.Status == 1){
                        GetShareFiles()
                    }else {
                        ErrMsg(resp.data.Msg)
                    }
                }).catch(err=>{
                    ErrMsg(err)
            })
        }



    },
    created: function () {
        this.mainMenuIndex=1
        window.addEventListener('resize', this.getHeight);
        this.getHeight()
        GetFiles(0)


    },
    destroyed:function () {
        window.removeEventListener('resize', this.getHeight);
    }
})

// file download
function download(content,fileName){
    const blob = new Blob([content]) //创建一个类文件对象：Blob对象表示一个不可变的、原始数据的类文件对象
    const url = window.URL.createObjectURL(blob)//URL.createObjectURL(object)表示生成一个File对象或Blob对象
    let dom = document.createElement('a')//设置一个隐藏的a标签，href为输出流，设置download
    dom.style.display = 'none'
    dom.href = url
    dom.setAttribute('download',fileName)//指示浏览器下载url,而不是导航到它；因此将提示用户将其保存为本地文件
    document.body.appendChild(dom)
    dom.click()
}

// get the use file by p
function GetFiles(p) {
    p = p||(vm?vm.$data.parentDir:0)||0
    axios.get('/file/userindexfiles?p='+p).then(function (response) {
        if (response.data.Status == 1){
            if (!response.data.Data){
                vm.tableData=[]
                return
            }

            vm.tableData = response.data.Data
        }else {
            alert(response.data.Msg)
        }
    }).catch(function (err) {
        console.log(err)
    })
}

// get the share file
function GetShareFiles() {
    vm.$data.appHost = window.location.href+"share/"
    axios.get("/file/usersharefiles/"+username)
        .then(resp=>{
            if (resp.data.Status ==1){
                vm.$data.shareTableData = resp.data.Data
            }else {
                ErrMsg(err)
            }
        })
        .catch(err=>{
            ErrMsg(err)
        })

}


function RightMenuDisplayNone() {
    document.querySelector("#context-menu").style.display = 'none';
}

function ErrMsg(msg) {
    vm.$message({
        message: msg,
        type: 'error',
        offset: 100
    });
}

//ref :https://blog.csdn.net/qq_34292479/article/details/87621491
function randomCode(){
    var arr = ['A','B','C','D','E','F','G','H','I','J','K','L','M','N','O','P','Q','R','S','T','U','V','W','X','Y','Z','a','b','c','d','e','f','g','h','i','j','k','l','m','n','o','p','q','r','s','t','u','v','w','x','y','z','0','1','2','3','4','5','6','7','8','9'];
    var idvalue ='';
    var n = 4;
    for(var i=0;i<n;i++){
        idvalue+=arr[Math.floor(Math.random()*62)];
    }
    return idvalue;
}

//ref: https://www.jb51.net/article/53061.htm
function string62to10(number_code) {
    var chars = '0123456789abcdefghigklmnopqrstuvwxyzABCDEFGHIGKLMNOPQRSTUVWXYZ',
        radix = chars.length,
        number_code = String(number_code),
        len = number_code.length,
        i = 0,
        origin_number = 0;
    while (i < len) {
        origin_number += Math.pow(radix, i++) * chars.indexOf(number_code.charAt(len - i) || 0);
    }
    return origin_number;
}



var flow = new Flow({
    target:'/file/upload',
    query:function (file, chunk, isTest) {
        return{
            'qetag':file.qetag
        }
    },
    chunkSize:ChunkSize,
    testChunks: false,
    headers:{'X-Requested-With':'XMLHttpRequest'}

});

// Flow.js isn't supported
if(!flow.support) {
    ErrMsg("上传插件不支持，请更换浏览器")


}

flow.assignBrowse(document.getElementById('browseButton'));


flow.on('fileAdded', function(file, event){
    if (vm.$data.drawer !=true){
        vm.$data.drawer =true
    }

    file.curprogress = 0
    file.uploadstatus = ''
    file.isSecond = 0
    vm.uploadTableData.push(file)
    get_qetag_and_upload(file)
});
flow.on('filesSubmitted', function(file, event){

});

flow.on('fileSuccess', function(file,message,chunk){
    //console.log("upload succ:"+file)
    var bodyFormData =new FormData()
    bodyFormData.append("qetag",file.qetag)
    bodyFormData.append("flowIdentifier",file.uniqueIdentifier)
    bodyFormData.append("flowTotalChunks",file.chunks.length)
    bodyFormData.append("fileSize",file.size)
    bodyFormData.append("fileName",file.name)
    bodyFormData.append("fileExt",file.name.split('.')[file.name.split('.').length-1])
    bodyFormData.append("parentDir",vm.$data.parentDir)
    axios({
        method: 'post',
        url: '/file/uploadfinshed',
        data: bodyFormData,
        headers: {'Content-Type': 'multipart/form-data'}
    }). then(response=> {
        console.log(response.data)
        if (parseInt(response.data.Status) ==1){
            file.uploadstatus="success"
            GetFiles()
        }else {
            file.uploadstatus="exception"
        }
    })


});

flow.on('fileProgress',function (file, chunk) {
    file.curprogress =Math.ceil( file.progress()*100)
})

flow.on('fileError', function(file, message,chunk){
    if (401 === chunk.xhr.status) {
        window.location = '/login';
    }else {
        ErrMsg("error:"+chunk.xhr.status)

    }

});



//qetag offical:https://github.com/qiniu/qetag
//get_qetag_and_upload ,custom function to get qetag ,then upload file
//custom the qetag function by flow.js
// sha1算法
var shA1 =  sha1.digest;
function get_qetag_and_upload(file) {
    var prefix = 0x16;
    var sha1String = [];

    var blockSize = ChunkSize;
    var blockCount = 0;
    var chunkIndex = 0


    function SetSha1StringArry(chunk) {

        var fileObj = chunk.fileObj
        var blob = webAPIFileRead(fileObj,chunk.startByte,chunk.endByte,fileObj.file.type)
        blob.arrayBuffer().then(function (buffer) {
            if (buffer == blockSize){
                sha1String.push(shA1(buffer));
            }else {
                var bufferSize = buffer.size || buffer.length || buffer.byteLength;
                blockCount = Math.ceil(bufferSize / blockSize);
                for (var i = 0; i < blockCount; i++) {
                    //sha1beginTime = +new Date();
                    sha1String.push(shA1(buffer.slice(i * blockSize, (i + 1) * blockSize)));
                    //var sha1endTime = +new Date();
                    //sha1CostTime +=sha1endTime-sha1beginTime


                }
            }

            chunkIndex++
            if (chunkIndex < file.chunks.length){
                //Forced synchronous execution to prevent qetag errors
                //TODO: can be optimized
                SetSha1StringArry(file.chunks[chunkIndex])

            }else {
                //to calc
                calcEtag(sha1String,file)
                //console.log("sha1共用时"+sha1CostTime+"ms");
            }
        })

    }
    // copy and modify by flow.js
    function webAPIFileRead(fileObj, startByte, endByte, fileType) {
        var function_name = 'slice';

        if (fileObj.file.slice)
            function_name =  'slice';
        else if (fileObj.file.mozSlice)
            function_name = 'mozSlice';
        else if (fileObj.file.webkitSlice)
            function_name = 'webkitSlice';

        return (fileObj.file[function_name](startByte, endByte, fileType));
    }

    SetSha1StringArry(file.chunks[chunkIndex])
    // for (var i=0; i< file.chunks.length ;i++){
    //     SetSha1StringArry(file.chunks[i])
    // }
    //return (calcEtag());
}

function calcEtag(sha1,file) {
    var prefix = 0x16;
    var sha1String = sha1;


    function concatArr2Uint8(s) {//Array 2 Uint8Array
        var tmp = [];
        for (var i of s) tmp = tmp.concat(i);
        return new Uint8Array(tmp);
    }
    function Uint8ToBase64(u8Arr, urisafe) {//Uint8Array 2 Base64
        var CHUNK_SIZE = 0x8000; //arbitrary number
        var index = 0;
        var length = u8Arr.length;
        var result = '';
        var slice;
        while (index < length) {
            slice = u8Arr.subarray(index, Math.min(index + CHUNK_SIZE, length));
            result += String.fromCharCode.apply(null, slice);
            index += CHUNK_SIZE;
        }
        return urisafe ? btoa(result).replace(/\//g, '_').replace(/\+/g, '-') : btoa(result);
    }
    function calcEtag() {
        if (!sha1String.length) return 'Fto5o-5ea0sNMlW_75VgGJCv2AcJ';
        var sha1Buffer = concatArr2Uint8(sha1String);
        // 如果大于4M，则对各个块的sha1结果再次sha1
        // ChunkSize = 4M
        if (file.size > ChunkSize) {
            prefix = 0x96;
            sha1Buffer = shA1(sha1Buffer.buffer);
        } else {
            sha1Buffer = Array.apply([], sha1Buffer);
        }
        sha1Buffer = concatArr2Uint8([[prefix], sha1Buffer]);
        var _base64 =Uint8ToBase64(sha1Buffer, true);

        //console.log("_base64 :"+_base64)
        return _base64
    }
    var qetag = calcEtag()
    file.qetag = qetag
    uploadfile(file)
}

function uploadfile(file) {

    fileSecondsPass(function (succ) {
        if (!succ){
            flow.upload()
        }else {
            //the seconds pass succ
            file.isSecond = 1
            file.cancel()
            GetFiles(vm.$data.parentDir)
        }
    },file)

}
//fileSecondsPass if can seconds-pass return true,else false
function fileSecondsPass(f,file) {
    var bodyFormData =new FormData()
    bodyFormData.append("qetag",file.qetag)
    bodyFormData.append("fileName",file.name)
    bodyFormData.append("parentDir",vm.$data.parentDir)
    axios.post("/file/filesecondspass",bodyFormData)
        .then(resp=>{
            if (parseInt(resp.data.Status)==1){
                f(true)
            }else {
                f(false)
            }
        }).catch(err=>{
            f(false)
    })

}


