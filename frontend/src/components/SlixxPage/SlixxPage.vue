<template>
  <div>
    <div v-if="totalPages <= 6">
      <ul class="pagination justify-content-end">
        <li :class="'page-item page-link ' + (page !== 1 ? '' : 'disabled-page-link')"
            @click="previousPage">
          Previous
        </li>
        <li :class="'page-item page-link ' + (page === target ? 'active-page-item' : '')"
            v-for="target in pages()"
            @click="onChange(target)">
          {{ target }}
        </li>
        <li :class="'page-item page-link ' + (page !== totalPages ? '' : 'disabled-page-link')"
            @click="nextPage">
          Next
        </li>
      </ul>
    </div>
    <div v-else>
      <select v-model="pageSelect" class="form-control pagination-select" @change="onChange(pageSelect * 1)">
        <option v-for="target in pages()" :value="target">Page {{ target }}</option>
      </select>
    </div>
  </div>
</template>
<script>

export default {
  name: 'Loader',
  props: {
    totalPages: {type: Number, default: 1},
    page: {type: Number, default: 1},
    onChange: {
      type: Function, default: (page) => {
      }
    }
  },
  data() {
    return {
      pageSelect: 1
    }
  },
  methods: {
    pages() {
      let pages = []
      for (let i = 1; i <= this.totalPages; i++) {
        pages.push(i)
      }
      return pages
    },
    previousPage() {
      if (this.page !== 1) {
        this.onChange(this.page - 1)
      }
    },
    nextPage() {
      if (this.page !== this.totalPages) {
        this.onChange(this.page + 1)
      }
    }
  }
}
</script>

<style src="./SlixxPage.scss" lang="scss"/>
