<template>
<el-card class="box-card model-card">
    <div slot="header" class="clearfix" v-drag>
        <span class="card-name">{{model.title}}</span>
        <el-link class="float-right" @click="$emit('deleteModel', model.id)">
            <i class="el-icon-close"/>
        </el-link>
    </div>
    <div class="text item">
        tttttt
    </div>
    <el-link class="resize-button" type="primary" v-resize>
        <i class="iconfont icon--Resize-Four-Direc"/>
    </el-link>
</el-card>
</template>

<script>
import $ from "jquery"

export default {
    props: {
        "model": Object
    },
    directives: {
        drag: { bind(el) {
            el.onmousedown = me => {
                let card = $(el).closest(".model-card")[0]
                let left = Number(card.style.left.slice(0, -2))
                let top = Number(card.style.top.slice(0, -2))
                let downClientX = me.clientX
                let downClientY = me.clientY

                document.onmousemove = e => {
                    card.style.left = `${left + (e.clientX - downClientX)}px`
                    card.style.top = `${top + (e.clientY - downClientY)}px`
                }
                document.onmouseup = e => {
                    document.onmousemove = null
                    document.onmouseup = null
                }
            }
        }},
        resize: { bind(el) {
            el.onmousedown = me => {
                let card = $(el).closest(".model-card")[0]
                let width = Number(card.style.width.slice(0, -2))
                let height = Number(card.style.height.slice(0, -2))
                let downClientX = me.clientX
                let downClientY = me.clientY

                document.onmousemove = e => {
                    card.style.width = `${e.clientX - downClientX}px`
                    card.style.height = `${e.clientY - downClientY}px`
                }
                document.onmouseup = e => {
                    document.onmousemove = null
                    document.onmouseup = null
                }
            }
        }}
    }
}
</script>


<style lang="scss">
.model-card {
    position: relative;
    width: 30vw;
    -webkit-touch-callout: none;
    -webkit-user-select: none;
    -khtml-user-select: none;
    -moz-user-select: none;
    -ms-user-select: none;
    user-select: none;

    .el-card__header {
        cursor: pointer;
    }
}
.resize-button {
    float: right !important;
    margin-bottom: 20px;
}
</style>
