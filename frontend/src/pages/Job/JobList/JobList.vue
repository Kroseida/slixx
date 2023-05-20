<template>
  <div class="tables-basic">
    <h2 class="page-title">Job List</h2>
    <b-row>
      <b-col>
        <Widget>
          <div class="table-pre-header">
            <b-button type="submit"
                      variant="outline-primary"
                      class="table-action"
                      id="action__csv__export"
                      @click="downloadCsv">
              <i class="fa fa-download px-1"></i>
              Export to CSV
            </b-button>
            <b-button type="submit"
                      variant="outline-success"
                      class="table-action"
                      v-if="isPermitted('job.create')"
                      id="action__create__job"
                      @click="openNew">
              <i class="fa fa-save px-1"></i>
              Create Job
            </b-button>
            <input class="form-control no-border table-search"
                   type="text"
                   name="search"
                   placeholder="Search"
                   v-model="table.search"
                   id="job__search"
                   v-on:keyup.enter="onSearch"
                   required/>
          </div>
          <SlixxTable :headers="headers" :entries="jobs" :onRowClick="open">
            <template v-slot:createdAt="slotProps">
              <div>
                {{ slotProps.entry.createdAt | moment('MM-DD-YYYY') }}
              </div>
            </template>
            <template v-slot:updatedAt="slotProps">
              <div>
                {{ slotProps.entry.updatedAt | moment('MM-DD-YYYY') }}
              </div>
            </template>
          </SlixxTable>
          <div class="table-post-item">
            <label class="table-item-count">
              Entries {{(table.page - 1) * 10}} to {{ (table.page - 1) * 10 + jobs.length }} (total: {{ table.totalRows }})
            </label>
            <SlixxPage :total-pages="table.totalPages" :on-change="onPageChange" :page="table.page"/>
          </div>
        </Widget>
      </b-col>
    </b-row>
  </div>
</template>

<script>
import SlixxTable from "@/components/SlixxTable/SlixxTable.vue";
import SlixxPage from "@/components/SlixxPage/SlixxPage.vue";
import {mapState, mapActions} from 'vuex';

export default {
  name: 'JobList',
  components: {
    SlixxTable,
    SlixxPage
  },
  data() {
    return {
      headers: [
        {
          name: 'Name',
          field: 'name'
        },
        {
          name: 'Strategy',
          field: 'strategy'
        },
        {
          name: 'Created At',
          field: 'createdAt'
        },
        {
          name: 'Updated At',
          field: 'updatedAt'
        }
      ],
    }
  },
  computed: {
    ...mapState('jobs', ['jobs', 'table']),
    ...mapState('layout', ['localUser', 'isPermitted'])
  },
  methods: {
    ...mapActions('users', ['subscribeJobs', 'unsubscribeJobs']),
    open(user) {
      this.$router.push({name: 'JobDetails', params: {id: user.id}});
    },
    openNew() {
      this.$router.push({name: 'JobDetails', params: {id: 'new'}});
    },
    downloadCsv() {
      let csv = '';
      csv += this.headers.map(header => header.field).join(';') + '\r';

      this.jobs.forEach(job => {
        csv += this.headers.map(header => job[header.field]).join(';') + '\r';
      });

      let blob = new Blob([csv], {type: 'text/csv;charset=utf-8;'});
      let link = document.createElement("a");
      let url = URL.createObjectURL(blob);
      link.setAttribute("href", url);
      link.setAttribute("download", 'jobs.csv');
      link.style.visibility = 'hidden';
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
    },
    onSearch() {
      this.$store.commit('jobs/unsubscribeJobs', {});
      this.table.page = 1;
      this.subscribeJobs();
    },
    subscribeJobs() {
      this.$store.commit('jobs/subscribeJobs', {
        error: (data) => {
          this.$toasted.error('Error while loading jobs: ' + data.message, {
            duration: 5000,
            position: 'top-right',
            fullWidth: true,
            fitToScreen: true,
          });
        }
      });
    },
    onPageChange(page) {
      this.$store.commit('jobs/unsubscribeJobs', {});
      this.table.page = page;
      this.subscribeJobs();
    }
  },
  mounted() {
    this.$store.commit('jobs/reset');
    this.subscribeJobs();
  },
  destroyed() {
    this.$store.commit('jobs/unsubscribeJobs', {});
    this.$store.commit('jobs/reset');
  }
};
</script>

<style src="./JobList.scss" lang="scss" scoped/>
