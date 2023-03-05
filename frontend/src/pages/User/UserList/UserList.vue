<template>
  <div>
    <h2 class="page-title">User List</h2>
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
                      v-if="isPermitted('user.create')"
                      id="action__create__user"
                      @click="openNew">
              <i class="fa fa-save px-1"></i>
              Create User
            </b-button>
            <input class="form-control no-border table-search"
                   type="text"
                   name="search"
                   placeholder="Search"
                   v-model="table.search"
                   id="user__search"
                   v-on:keyup.enter="onSearch"
                   required/>
          </div>
          <SlixxTable :headers="headers" :entries="users" :onRowClick="open">
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
            <template v-slot:active="slotProps">
              <div>
                <b-badge v-if="slotProps.entry.active" variant="success" pill>
                  Active
                </b-badge>
                <b-badge v-else variant="danger" pill>
                  Inactive
                </b-badge>
              </div>
            </template>
          </SlixxTable>
          <div class="table-post-item">
            <label class="table-item-count">
              Entries {{(table.page - 1) * 10}} to {{ (table.page - 1) * 10 + users.length }} (total: {{ table.totalRows }})
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
  name: 'UserList',
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
          name: 'First Name',
          field: 'firstName'
        },
        {
          name: 'Last Name',
          field: 'lastName'
        },
        {
          name: 'Created At',
          field: 'createdAt'
        },
        {
          name: 'Updated At',
          field: 'updatedAt'
        },
        {
          name: 'Status',
          field: 'active'
        }
      ],
    }
  },
  computed: {
    ...mapState('users', ['users', 'table']),
    ...mapState('layout', ['localUser', 'isPermitted'])
  },
  methods: {
    ...mapActions('users', ['subscribeUsers', 'unsubscribeUsers']),
    open(user) {
      this.$router.push({name: 'UserDetails', params: {id: user.id}});
    },
    openNew() {
      this.$router.push({name: 'UserDetails', params: {id: 'new'}});
    },
    downloadCsv() {
      let csv = '';
      csv += this.headers.map(header => header.field).join(';') + '\r';

      this.users.forEach(user => {
        csv += this.headers.map(header => user[header.field]).join(';') + '\r';
      });

      let blob = new Blob([csv], {type: 'text/csv;charset=utf-8;'});
      let link = document.createElement("a");
      let url = URL.createObjectURL(blob);
      link.setAttribute("href", url);
      link.setAttribute("download", 'users.csv');
      link.style.visibility = 'hidden';
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
    },
    onSearch() {
      this.$store.commit('users/unsubscribeUsers', {});
      this.table.page = 1;
      this.subscribeUsers();
    },
    subscribeUsers() {
      this.$store.commit('users/subscribeUsers', {
        error: (data) => {
          this.$toasted.error('Error while loading users: ' + data.message, {
            duration: 5000,
            position: 'top-right',
            fullWidth: true,
            fitToScreen: true,
          });
        }
      });
    },
    onPageChange(page) {
      this.$store.commit('users/unsubscribeUsers', {});
      this.table.page = page;
      this.subscribeUsers();
    }
  },
  mounted() {
    this.$store.commit('user/reset');
    this.subscribeUsers();
  },
  destroyed() {
    this.$store.commit('users/unsubscribeUsers', {});
    this.$store.commit('user/reset');
  }
};
</script>

<style src="./UserList.scss" lang="scss" scoped/>
