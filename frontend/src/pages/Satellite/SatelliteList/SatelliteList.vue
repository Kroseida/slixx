<template>
  <div class="tables-basic">
    <h2 class="page-title">Satellite List</h2>
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
                      v-if="isPermitted('satellite.create')"
                      id="action__create__satellite"
                      @click="openNew">
              <i class="fa fa-save px-1"></i>
              Register Satellite
            </b-button>
            <input class="form-control no-border table-search"
                   type="text"
                   name="search"
                   placeholder="Search"
                   v-model="table.search"
                   id="satellite__search"
                   v-on:keyup.enter="onSearch"
                   required/>
          </div>
          <SlixxTable :headers="headers" :entries="satellites" :onRowClick="open">
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
              Entries {{(table.page - 1) * 10}} to {{ (table.page - 1) * 10 + satellites.length }} (total: {{ table.totalRows }})
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
  name: 'SatelliteList',
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
          name: 'Address',
          field: 'address'
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
    ...mapState('satellites', ['satellites', 'table']),
    ...mapState('layout', ['localUser', 'isPermitted'])
  },
  methods: {
    open(satellite) {
      this.$router.push({name: 'SatelliteDetails', params: {id: satellite.id}});
    },
    openNew() {
      this.$router.push({name: 'SatelliteDetails', params: {id: 'new'}});
    },
    downloadCsv() {
      let csv = '';
      csv += this.headers.map(header => header.field).join(';') + '\r';

      this.satellites.forEach(satellite => {
        csv += this.headers.map(header => satellite[header.field]).join(';') + '\r';
      });

      let blob = new Blob([csv], {type: 'text/csv;charset=utf-8;'});
      let link = document.createElement("a");
      let url = URL.createObjectURL(blob);
      link.setAttribute("href", url);
      link.setAttribute("download", 'satellites.csv');
      link.style.visibility = 'hidden';
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
    },
    onSearch() {
      this.$store.commit('satellites/unsubscribeSatellites', {});
      this.table.page = 1;
      this.subscribeSatellites();
    },
    subscribeSatellites() {
      this.$store.commit('satellites/subscribeSatellites', {
        error: (data) => {
          this.$toasted.error('Error while loading satellites: ' + data.message, {
            duration: 5000,
            position: 'top-right',
            fullWidth: true,
            fitToScreen: true,
          });
        }
      });
    },
    onPageChange(page) {
      this.$store.commit('satellites/unsubscribeSatellites', {});
      this.table.page = page;
      this.subscribeSatellites();
    }
  },
  mounted() {
    this.$store.commit('satellites/reset');
    this.subscribeSatellites();
  },
  destroyed() {
    this.$store.commit('satellites/unsubscribeSatellites', {});
    this.$store.commit('satellites/reset');
  }
};
</script>

<style src="./SatelliteList.scss" lang="scss" scoped/>
