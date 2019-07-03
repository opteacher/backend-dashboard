<template>
    <div class="w-100 h-100">
        <tool-bar @add-model="addModel" @add-link="addLink" :models="models"/>
        <svg id="pnlModels" class="w-100 h-100"></svg>
    </div>
</template>

<script>
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
            this.queryModels()
            // this.queryLinks()
        },
        watch: {
            models() {
                d3.select("#pnlModels")
                    .selectAll("circle")
                    .data(this.models)
                    .enter()
                    .append("circle")
                    .attr("cx", model => model.x)
                    .attr("cy", model => model.y)
                    .attr("r", 20)
                    .attr("fill", "steelblue")
                    .call(d3.drag().on("drag", function (tgt) {
                        d3.select(this)
                            .attr("cx", tgt.cx = d3.event.x)
                            .attr("cy", tgt.cx = d3.event.y)
                    }))
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
