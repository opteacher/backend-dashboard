<template>
    <div class="w-100 h-100">
        <tool-bar @add-model="addModel" @add-relation="addRelation" :models="models"/>
        <div id="pnlModels" class="w-100 h-100"></div>
    </div>
</template>

<script>
    import echarts from "echarts"
    import toolBar from "../components/toolBar"
    import modelBkd from "../async/model"
    import relationBkd from "../async/relation"

    export default {
        components: {
            "tool-bar": toolBar,
        },
        data() {
            return {
                models: [],
                relations: []
            }
        },
        created() {
            this.queryModels()
            this.queryRelations()
        },
        watch: {
            models() {
                console.log(this.models)
                let graph = {  //这是数据项目中一般都是获取到的
                    nodes: this.models,
                    links: []
                }
                let myChart = echarts.init(document.getElementById('pnlModels'))
                graph.nodes.forEach(node => {
                    node.x = parseInt(Math.random() * 1000)  //这里是最重要的如果数据中有返回节点x,y位置这里就不用设置，如果没有这里一定要设置node.x和node.y，不然无法定位节点 也实现不了拖拽了；
                    node.y = parseInt(Math.random() * 1000)
                })
                let option = {    //这里是option配置
                    animationDurationUpdate: 1500,
                    animationEasingUpdate: 'quinticInOut',
                    series: [
                        {
                            type: 'graph',
                            layout: 'none',           //因为节点的位置已经有了就不用在这里使用布局了
                            symbolSize: 50,
                            circular: {rotateLabel: true},
                            animation: false,
                            data: graph.nodes,
                            links: graph.links,
                            roam: true,   //添加缩放和移动
                            draggable: false,   //注意这里设置为false，不然拖拽鼠标和节点有偏移
                            label: {
                                normal: {
                                    show: true
                                }
                            }
                        }
                    ]
                }
                myChart.setOption(option)
                initInvisibleGraphic()

                function initInvisibleGraphic() {
                    // Add shadow circles (which is not visible) to enable drag.
                    myChart.setOption({
                        graphic: echarts.util.map(option.series[0].data, (item, dataIndex) => {
                            //使用图形元素组件在节点上划出一个隐形的图形覆盖住节点
                            let tmpPos = myChart.convertToPixel({'seriesIndex': 0}, [item.x, item.y])
                            return {
                                type: 'circle',
                                id: dataIndex,
                                position: tmpPos,
                                shape: {
                                    cx: 0,
                                    cy: 0,
                                    r: 50
                                },
                                // silent:true,
                                invisible: true,
                                draggable: true,
                                ondrag: echarts.util.curry(onPointDragging, dataIndex),
                                z: 100              //使图层在最高层
                            }
                        })
                    })
                    window.addEventListener('resize', updatePosition)
                    myChart.on('dataZoom', updatePosition)
                }

                myChart.on('graphRoam', updatePosition)

                function updatePosition() {    //更新节点定位的函数
                    myChart.setOption({
                        graphic: echarts.util.map(option.series[0].data, item => {
                            position: myChart.convertToPixel({'seriesIndex': 0}, [item.x, item.y])
                        })
                    })

                }

                function onPointDragging(dataIndex) {      //节点上图层拖拽执行的函数
                    let tmpPos = myChart.convertFromPixel({'seriesIndex': 0}, this.position)
                    option.series[0].data[dataIndex].x = tmpPos[0]
                    option.series[0].data[dataIndex].y = tmpPos[1]
                    myChart.setOption(option)
                    updatePosition()
                }
            }
        },
        methods: {
            async queryModels() {
                let res = await modelBkd.qry()
                if (typeof res === "string") {
                    this.$message(`查询模块失败：${res}`)
                } else {
                    this.models = (res.data.data && res.data.data.models) || []
                }
            },
            async queryRelations() {
                let res = await relationBkd.qry()
                if (typeof res === "string") {
                    this.$message(`查询关联失败：${res}`)
                } else {
                    this.relations = res.data.data || []
                }
            },
            async addModel(model) {
                let res = await modelBkd.add(model)
                if (typeof res === "string") {
                    this.$message(`创建模块失败：${res}`)
                } else {
                    this.models.push(model)
                }
            },
            async deleteModel(modelID) {
                let res = await modelBkd.del(modelID)
                if (typeof res === "string") {
                    this.$message(`删除模块失败：${res}`)
                } else {
                    this.models.pop(ele => ele.id === modelID)
                }
            },
            async addRelation(relation) {
                let res = await relationBkd.add(relation)
                if (typeof res === "string") {
                    this.$message(`创建关联失败：${res}`)
                } else {
                    relation.id = res.data.data[0].id
                    this.relations.push(relation)
                }
            }
        },
    }
</script>
