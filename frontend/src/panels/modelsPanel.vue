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
            chart: null,
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
            this.chart = echarts.init(document.getElementById("pnlModels"))
            this.chart.setOption({
                title: {},
                tooltip: {},
                animationDurationUpdate: 1500,
                animationEasingUpdate: 'quinticInOut',
                series : [
                    {
                        type: 'graph',
                        layout: 'none',
                        symbolSize: 50,
                        roam: true,
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
                        data: this.models,
                        links: [],
                        // links: [{
                        //     source: 0,
                        //     target: 1,
                        //     symbolSize: [5, 20],
                        //     label: {
                        //         normal: {
                        //             show: true
                        //         }
                        //     },
                        //     lineStyle: {
                        //         normal: {
                        //             width: 5,
                        //             curveness: 0.2
                        //         }
                        //     }
                        // }, {
                        //     source: '节点2',
                        //     target: '节点1',
                        //     label: {
                        //         normal: {
                        //             show: true
                        //         }
                        //     },
                        //     lineStyle: {
                        //         normal: { curveness: 0.2 }
                        //     }
                        // }, {
                        //     source: '节点1',
                        //     target: '节点3'
                        // }, {
                        //     source: '节点2',
                        //     target: '节点3'
                        // }, {
                        //     source: '节点2',
                        //     target: '节点4'
                        // }, {
                        //     source: '节点1',
                        //     target: '节点4'
                        // }],
                        // lineStyle: {
                        //     normal: {
                        //         opacity: 0.9,
                        //         width: 2,
                        //         curveness: 0
                        //     }
                        // }
                    }
                ]
            });
        }
    },
    methods: {
        async queryModels() {
            let res = await modelBkd.qry()
            if (typeof res === "string") {
                this.$message(`查询模块失败：${res}`)
            } else {
                this.models = res.data.data.models || []
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
                model.id = res.data.data[0].id
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
        },
        findModel(id) {
            return this.models.find(ele => ele.id === id)
        },
        bindRelationToModel(modelID, relation) {
            let model = this.findModel(modelID)
            model.observers = [relation]
            model.notifyUpdate = function() {
                for (let obs of this.observers) {
                    obs.onModelChanged(this.x, this.y, this.width, this.height)
                }
            }
        }
    },
}
</script>
