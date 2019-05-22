<template>
<el-card class="box-card model-card" style="width: 300px">
    <div slot="header" class="clearfix" v-drag>
        <span class="card-name">{{model.name}}</span>
        <el-link class="float-right" @click="$emit('deleteModel', model.id)">
            <i class="el-icon-close"/>
        </el-link>
    </div>
    <div class="text item">
        <div class="list-group list-group-flush">
            <button type="button" class="list-group-item list-group-item-action list-group-item-light" v-for="prop in model.props" :key="prop.name">{{prop.name}}</button>
        </div>
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
                    let l = left + (e.clientX - downClientX)
                    let t = top + (e.clientY - downClientY)
                    t = (t <= 0) ? 0 : t
                    l = (l <= 0) ? 0 : l
                    card.style.left = `${l}px`
                    card.style.top = `${t}px`
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
                    let w = width +  e.clientX - downClientX
                    let h = height + e.clientY - downClientY
                    h = (h <= 61 + 24 + 40) ? 61 + 24 + 40 : h
                    card.style.width = `${w}px`
                    card.style.height = `${h}px`
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
    position: absolute;
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
    position: absolute;
    right: 20px;
    bottom: 20px;
}
</style>
