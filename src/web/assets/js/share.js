
Vue.prototype.$axios = axios
var vm =new Vue({
    delimiters: ['$', '$'],//use the '$$' delimiters,because the '{{' confilct with 'iris' mvc
    el: '#app',
    data:{
        file_status:0,
        share_pwd:'',
        need_pwd:1,
        share_id:'',
        show_content:0,
        is_dir:0,
        share_time:"",
        create_time:"",
        file_name:"",
        cur_file:{},
        moveFileDialogFormVisible:false,
        moveFileData:[],
        moveFileTreeDefaultExpanded:[],
        moveFiledefaultProps: {
            children: 'children',
            label: 'label'
        },

    },
    methods: {
        was_login(){
            return username.length >0
        },
        onValid(){
            form = new FormData()
            form.append("share_id",this.share_id)
            form.append("pwd",this.share_pwd)

            this.$axios.post("/share/valid",form)
                .then(resp=>{
                    if (resp.data.Status==0){
                        ErrMsg(resp.data.Msg)
                        return
                    }
                    this.show_content = 1
                    if(resp.data.Data.IsDir==0){
                        //is file,show
                        this.is_dir = 0
                        this.create_time = resp.data.Data.CreateAt
                        this.file_name = resp.data.Data.FileName
                        this.cur_file = resp.data.Data
                        if (resp.data.Data.ShareTime==1){
                            this.share_time = "1 天"

                        }else if(resp.data.Data.ShareTime==7){
                            this.share_time = "7 天"

                        }else {
                            this.share_time = "永久"
                        }

                    }else {
                        //is dir
                        //to get the dir content
                        this.is_dir = 1

                    }

                }).catch(err=>{
                    ErrMsg(err)
            })
        },
        mouseleftclick(){
            var menu = document.querySelector("#context-menu");
            menu.style.display = 'none';
            this.curRightRow = {}
        },
        getHeight(){
            this.tableheight=window.innerHeight-121.511+'px';  //获取浏览器高度减去顶部导航栏
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
            _obj  = row.LastUpdated.toString().split(':')
            return  row.LastUpdated.toString().replace(":"+_obj[_obj.length-1],"")
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
        allfile(index,indexPath){

            if(index==1){
                this.NavArray =[{
                    "ID":0,
                    "FileName":"首页"
                }]
                this.parentDir = 0
                GetFiles(0)
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

            if (this.cur_file.IsDir==1){
                ErrMsg("暂不支持文件夹下载功能")
                return
            }
            try {
                var elemIF = document.createElement("iframe");
                elemIF.src = "/share/downloadfile/"+this.cur_file.FileName+"?share_id="+this.cur_file.ShareId+"&share_pwd="+this.share_pwd;
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

            this.$axios.get("/file/userdirs/"+this.cur_file.UserFileId).
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
            data.append("share_id",this.share_id)
            data.append("share_pwd",this.share_pwd)
            data.append("dir",this.curMoveFileTreeSelected)

            this.$axios.post("/share/savefile",data).
            then(resp=>{
                if (resp.data.Status ==1){
                    vm.$message({
                        message: "文件保存成功",
                        type: 'success',
                        offset: 100
                    });
                }else {
                    ErrMsg(resp.data.Msg)
                }
            })
        }
    },
    created: function () {
        window.addEventListener('resize', this.getHeight);
        this.getHeight()
        this.share_id = document.querySelector("#share_id").value
        if (this.share_id.length==0){
            return
        }
        this.$axios.get("/share/file/"+this.share_id).
            then(resp=>{
                if (resp.data.Status ==1){
                    this.file_status = 1
                    if(resp.data.Data.pwd==1){
                        this.need_pwd = 1
                    }else {
                        this.need_pwd = 0
                    }
                }else{
                    //file status err
                    this.file_status = 0
            }
        }).catch(err=>{
            ErrMsg(err)
        })


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
