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
                let mdlCard = d3.select("#pnlModels")
                    .selectAll("g")
                    .data(this.models)
                    .join("g")
                let padding = 5
                let innerHgt = {}
                mdlCard.append("rect")
                    .attr("name", model => `model_${model.name}`)
                    .attr("x", model => model.x)
                    .attr("y", model => model.y)
                    .attr("width", model => model.width)
                    .attr("height", model => model.height)
                    .attr("rx", 4)
                    .attr("ry", 4)
                    .attr("fill", "steelblue")
                mdlCard.append("text")
                    .attr("x", model => model.x)
                    .attr("y", model => model.y)
                    .text(model => model.name)
                    .attr("fill", "white")
                    .each(function (model) {
                        innerHgt[model.name] = this.getBoundingClientRect().height
                        d3.select(this)
                            .attr("x", model.x + padding)
                            .attr("y", model.y + this.getBoundingClientRect().height)
                    })
                mdlCard.append("line")
                    .attr("x1", model => model.x + padding)
                    .attr("y1", model => model.y + innerHgt[model.name] + padding)
                    .attr("x2", model => model.x + model.width - padding)
                    .attr("y2", model => model.y + innerHgt[model.name] + padding)
                    .attr("stroke-width", 2)
                    .attr("stroke", "white")
                mdlCard.append("rect")
                    .attr("name", model => `model_${model.name}`)
                    .attr("x", model => model.x)
                    .attr("y", model => model.y)
                    .attr("width", model => model.width)
                    .attr("height", model => model.height)
                    .attr("rx", 4)
                    .attr("ry", 4)
                    .attr("opacity", 0)
                    .call(d3.drag().on("drag", function (tgt) {
                        d3.selectAll(`[name='model_${tgt.name}']`)
                            .attr("x", tgt.x = d3.event.x)
                            .attr("y", tgt.y = d3.event.y)
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
