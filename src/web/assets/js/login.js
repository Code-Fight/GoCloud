var vm =new Vue({
    el:"#app",
    data:{
        dowebok:'dowebok',
        pan:''
    },
    methods: {
        signInButton(){
            this.pan='right-panel-active'
        },
        signUpButton(){
            this.pan=''
        },
        signIn() {
            var user_name = document.getElementById('signup_user_name').value
            var eamil = document.getElementById('signupeamil').value
            var password = document.getElementById('signuppassword').value



            var data= {"Username":user_name,"Email":eamil,"Userpwd":password}
            axios.post('/login/signup',data).
            then(function (response) {
                if (response.data == 'OK'){

                    vm.$message({
                        message: '恭喜你，注册成功',
                        type: 'success'
                    });
                    vm.signUpButton()
                }else {
                    vm.$message({
                        message: '注册错误:'+response.data,
                        type: 'error'
                    });
                }
            }).catch(function (error) {
                console.log(error)

            })
        },
        signUp() {

            var eamil = document.getElementById('eamil').value
            var password = document.getElementById('passwd').value

            if (eamil.length ==0 || password.length ==0){
                vm.$message({
                    message: '请输入邮箱和密码',
                    type: 'error'
                });
                return
            }

            var data= {"Email":eamil,"Userpwd":password}
            axios.post('/login/signin',data).
            then(function (response) {
                if (response.data.Msg == 'OK'){

                    vm.$message({
                        message: '恭喜你，登录成功',
                        type: 'success'
                    });
                    window.location.href= "/"
                }else {
                    vm.$message({
                        message: '登录失败:'+response.data.Msg,
                        type: 'error'
                    });
                }
            }).catch(function (error) {
                console.log(error)
            })
        }
    }
})
