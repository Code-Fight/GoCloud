//the frontend chunksize must equal backend chunksize
var ChunkSize = 4*1024*1024

var beginTime
var sha1CostTime =0
Vue.prototype.$axios = axios
var vm =new Vue({
    delimiters: ['$', '$'],//use the '$$' delimiters,because the '{{' confilct with 'iris' mvc
    el: '#app',
    data:{
        tableheight:"auto",
        tableData: [],
        multipleSelection: []

    },
    methods: {
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
                    Math.round( parseFloat(row.FileSize)/1024*1024*1024 * 10) / 10
                    return  Math.round( parseFloat(row.FileSize)/1024*1024*1024 * 10) / 10+'G'
                }else if (parseInt(row.FileSize) >= 1024*1024){
                    return Math.round( parseFloat(row.FileSize)/1024*1024 * 10) / 10+'M'
                }else if (parseInt(row.FileSize) >= 1024){
                    return Math.round( parseFloat(row.FileSize)/1024 * 10) / 10+'K'
                }else {
                    return row.FileSize+'b'
                }
            }
        },
        table_date_formatter(row, column){
            _obj  = row.LastUpdated.toString().split(':')
            return  row.LastUpdated.toString().replace(":"+_obj[_obj.length-1],"")
        }
    },
    created: function () {
        window.addEventListener('resize', this.getHeight);
        this.getHeight()


        // 请求后端 获取值
        this.$axios.get('/file/userindexfiles').then(function (response) {
            if (response.data.Status == 1){
                vm.tableData = response.data.Data
            }else {
                alert(response.data.Msg)
            }
        }).catch(function (err) {
            console.log(err)
        })
    },
    destroyed:function () {
        window.removeEventListener('resize', this.getHeight);
    }
})


var flow = new Flow({
    target:'/file/upload',
    query:function (file, chunk, isTest) {
        return{
            'qetag':file.qetag
        }
    },
    chunkSize:ChunkSize,
    testChunks: false

});
// Flow.js isn't supported
if(!flow.support) {
    vm.$message({
        message: '上传插件不支持，请更换浏览器',
        type: 'error'
    });

}

flow.assignBrowse(document.getElementById('browseButton'));
// flow.assignDrop(document.getElementById('dropTarget'));

flow.on('fileAdded', function(file, event){

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
    axios({
        method: 'post',
        url: '/file/uploadfinshed',
        data: bodyFormData,
        headers: {'Content-Type': 'multipart/form-data'}
    }). then(response=> {
        alert(response.data.Msg)
    })


});
flow.on('fileError', function(file, message){
    alert("UPLOAD ERROR:"+message)
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
    // var endTime = +new Date();
    // console.log("用时共计"+(endTime-beginTime)+"ms");
    // console.log("upload qetag:"+qetag)

    if (!fileSecondsPass(file.qetag)){
        flow.upload()
    }else {
        // markup and tip the file was uploaded success
        // TODO
    }

    //
}
//fileSecondsPass if can seconds-pass return true,else false
function fileSecondsPass() {
    return false
}


