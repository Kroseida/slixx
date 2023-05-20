<template>
  <div class="data-loader">
    <i class="la la-spinner la-spin" :style="{fontSize: size + 'px'}"></i>
  </div>
</template>
<script>

export default {
  name: 'Loader',
  props: {
    size: {type: Number, default: 21}
  }
}
</script>

<style src="./Loader.scss" lang="scss"/>
