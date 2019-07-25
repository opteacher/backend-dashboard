<template>
<div class="card" :name="`flow_${step.index}`" :style="`
    width:${flowWidth}px;margin-bottom:${step.isLast ? marginTB : 0}px
`">
    <div class="card-header text-center">
        <h5 class="mb-0 float-left">#{{step.index}}</h5>
        <span>{{step.operKey}}</span>
        <button class="close" type="button" data-dismiss="alert" aria-label="Close" @click="hdlDelete">
            <span aria-hidden="true">&times;</span>
        </button>
    </div>
    <div class="row">
        <div class="col pr-0">
            <ul class="list-group list-group-flush h-100">
                <a class="list-group-item list-group-item-primary list-group-item-action api-params" href="#" v-for="(content, pholder) in step.inputs" :key="pholder">
                    {{pholder}}
                    <i class="el-icon-arrow-right"/>
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
</template>

<script>
export default {
    props: {
        "step": Object
    },
    data() {
        return {
            marginTB: 50,
            flowWidth: 600,
            flowSpan: 300,
        }
    },
    created() {
        let panelWidth = parseInt(document.getElementById("pnlFlows").getBoundingClientRect().width)
        this.step.x = (panelWidth>>1) - (this.flowWidth>>1)
        this.step.y = this.marginTB + this.flowSpan * this.step.index
    },
    mounted() {
        d3.select(`[name="flow_${this.step.index}"]`)
            .style("left", `${this.step.x}px`)
            .style("top", `${this.step.y}px`)
    },
    methods: {
        showOperDetail() {
            this.$emit("show-detail", this.step)
        },
        hdlDelete() {
            // TODO: 删除步骤
        }
    }
}
</script>

