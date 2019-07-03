<template>
    <div class="w-100 h-100">
        <tool-bar @add-model="addModel" @add-link="addLink" :models="models"/>
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
                mdlChart: null,
                models: [],
                links: []
            }
        },
        created() {
            // this.queryModels()
            // this.queryLinks()
            this.models = [{
                name: "User",
                x: 200,
                y: 200,
                fixed: true
            }, {
                name: "Company",
                x: 200,
                y: 300,
                fixed: true
            }]
            this.links = [{
                source: 0,
                target: 1
            }]
        },
        watch: {
            models() {
                if (!this.mdlChart) {
                    this.mdlChart = echarts.init(document.getElementById('pnlModels'))
                    let option = {
                        animationDurationUpdate: 1500,
                        animationEasingUpdate: 'quinticInOut',
                        series: [
                            {
                                type: 'graph',
                                layout: 'force',
                                symbolSize: 50,
                                roam: false,
                                draggable: false,
                                animation: false,
                                label: {
                                    normal: {
                                        show: true
                                    }
                                },
                                edgeSymbol: ['circle', 'arrow'],
                                edgeSymbolSize: [4, 10],
                                edgeLabel: {
                                    normal: {
                                        textStyle: {
                                            fontSize: 20
                                        }
                                    }
                                },
                            //     data: this.models,
                            //     links: this.links
                            }
                        ]
                    }
                    this.mdlChart.setOption(option)
                }
                let option = this.mdlChart.getOption()
                this.mdlChart.setOption({
                    graphic: echarts.util.map(this.models, (item, index) => {
                        return {
                            id: index,
                            type: "circle",
                            position: this.mdlChart.convertToPixel({'seriesIndex': 0}, [item.x, item.y]),
                            shape: {
                                cx: 0,
                                cy: 0,
                                r: 50/2
                            },
                            style: {
                                fill: "#c33531"
                            },
                            draggable: true,
                            ondrag: echarts.util.curry((dIdx, eve) => {
                                let pos = this.mdlChart.convertFromPixel({'seriesIndex': 0}, eve.target.position)
                                this.models[dIdx].x = pos[0]
                                this.models[dIdx].y = pos[1]
                                this.mdlChart.setOption({
                                    graphic: echarts.util.map(this.models, item => {
                                        position: this.mdlChart.convertToPixel({'seriesIndex': 0}, [item.x, item.y])
                                    })
                                })
                            }, index),
                            z: 100
                        }, {
                            type: "text",
                            style: {
                                "text": item.name
                            }
                        }
                    })
                })
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
            async queryLinks() {
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
            async addLink(relation) {
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
