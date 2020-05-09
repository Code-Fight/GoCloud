
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
        toggleSelection(rows) {
            if (rows) {
                rows.forEach(row => {
                    this.$refs.multipleTable.toggleRowSelection(row);
                });
            } else {
                this.$refs.multipleTable.clearSelection();
            }
        },
        handleSelectionChange(val) {
            this.multipleSelection = val;
        }
    }
})