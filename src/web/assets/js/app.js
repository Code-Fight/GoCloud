
Vue.prototype.$axios = axios
var vm =new Vue({
    el: '#app',
    data:{

        tableData: [],
        multipleSelection: []

    },
    created: function () {
        // 请求后端 获取值
        this.$axios.get('/index').then(function (response) {
            console.log(response)
        }).catch(function (err) {
            console.log(err)
        })
        // `this` 指向 vm 实例
        // this._data.tableData = [{
        //     date: '2016-05-03',
        //     name: '文件-文件-文件-文件',
        //     size: '2.3G'
        // }, {
        //     date: '2016-05-02',
        //     name: '文件-文件-文件-文件',
        //     size: '2.3G'
        // },{
        //     date: '2016-05-02',
        //     name: '文件-文件-文件-文件',
        //     size: '2.3G'
        // },{
        //     date: '2016-05-01',
        //     name: '文件-文件-文件-文件',
        //     size: '2.3G'
        // },{
        //     date: '2016-05-08',
        //     name: '文件-文件-文件-文件',
        //     size: '2.3G'
        // },{
        //     date: '2016-05-13',
        //     name: '文件-文件-文件-文件',
        //     size: '2.3G'
        // }]
    },


    methods: {
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