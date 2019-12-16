<template>
<div style="position:absolute;width:0;height:0">
    <div class="card" :name="`step_${step.index}`" :style="`width:${stepWidth}px;margin-bottom:${step.isLast ? marginTB : 0}px`">
        <div class="card-header text-center">
            <h5 class="mb-0 float-left">#{{step.index}}</h5>
            <span>{{step.key}}</span>
            <button class="close" type="button" data-dismiss="alert" aria-label="Close" @click="hdlDelete">
                <span aria-hidden="true">&times;</span>
            </button>
        </div>
        <div class="row">
            <div class="col pr-0">
                <ul class="list-group list-group-flush h-100">
                    <a class="list-group-item list-group-item-primary list-group-item-action api-params" href="#" v-for="(content, pholder) in step.inputs" :key="pholder">
                        {{content}}
                        <i class="el-icon-arrow-right"/>
                        {{pholder}}
                    </a>
                </ul>
            </div>
            <div class="col-6 card-body text-center desc-panel" @click="showOperDetail">{{step.desc}}</div>
            <div class="col pl-0">
                <div class="list-group list-group-flush h-100">
                    <a class="list-group-item list-group-item-success list-group-item-action api-params text-right" href="#" v-for="output in step.outputs" :key="output">
                        {{output}}
                        <i class="el-icon-arrow-right"/>
                    </a>
                </div>
            </div>
        </div>
    </div>
</div>
</template>

<script>
import backend from '../backend';
export default {
    props: {
        "step": Object
    },
    data() {
        return {
            marginTB: 50,
            stepWidth: 800,
            stepSpan: 300,
        }
    },
    mounted() {
        let panelWidth = parseInt(document.getElementById("pnlFlows").getBoundingClientRect().width)
        d3.select(`[name="step_${this.step.index}"]`)
            .style("left", `${(panelWidth>>1) - (this.stepWidth>>1)}px`)
            .style("top", `${this.marginTB + this.stepSpan * this.step.index}px`)
    },
    methods: {
        showOperDetail() {
            this.$emit("show-detail", this.step)
        },
        hdlDelete() {
            this.$alert("确定删除该步骤？", "提示", {
                confirmButtonText: "确定",
                callback: async action => {
                    if (action !== "confirm") {
                        return
                    }
                    let res = await backend.delStep({
                        apiName: this.step.apiName,
                        stepId: this.step.index
                    })
                    if (typeof res === "string") {
                        this.$message.error(`删除步骤时发生错误${res}`)
                    } else {
                        this.$message({
                            type: "info",
                            message: `步骤#${this.step.index} ${this.step.key}删除成功！`
                        })
                        this.$emit("be-deleted")
                    }
                }
            })
        }
    }
}
</script>

