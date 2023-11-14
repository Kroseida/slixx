import {useConstantsStore} from "stores/constants";

export default {
  name: "PermissionAssignComponent",
  components: {},
  emits: ['add-permission', 'remove-permission'],
  setup() {
    return {
      constantsStore: useConstantsStore(),
    };
  },
  data() {
    return {
      permissionMap: {},
    };
  },
  props: {
    currentPermissions: {
      type: Array,
      default: () => [],
    }
  },
  mounted() {
    this.constantsStore.permissions.forEach((permission) => {
      this.permissionMap[permission.value] = this.currentPermissions.includes(permission.value);
    });
  },
  methods: {
    togglePermission(permission, hasPermission) {
      if (hasPermission) {
        this.$emit('add-permission', permission)
      } else {
        this.$emit('remove-permission', permission)
      }
    },
    sortedPermissions() {
      return this.constantsStore.permissions.sort((a, b) => {
        return a.name.localeCompare(b.name);
      });
    }
  },
  watch: {
    currentPermissions: {
      handler: function (val) {
        this.constantsStore.permissions.forEach((permission) => {
          this.permissionMap[permission.value] = val.includes(permission.value);
        });
      }
    }
  }
};
