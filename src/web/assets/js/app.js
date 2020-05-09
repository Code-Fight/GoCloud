Vue.prototype.$axios = axios
var vm =new Vue({
    el: '#app',
    data:{
        tableheight:"",
        tableData: [],
        multipleSelection: []

    },
    created: function () {
        window.addEventListener('resize', this.getHeight);
        this.getHeight()


        // 请求后端 获取值
        this.$axios.get('/index').then(function (response) {
            console.log(response)
        }).catch(function (err) {
            console.log(err)
        })
    },
    destroyed:function () {
        window.removeEventListener('resize', this.getHeight);
    },

    methods: {
        getHeight(){
            this.tableheight=window.innerHeight-120+'px';  //获取浏览器高度减去顶部导航栏
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
        }
    }
})