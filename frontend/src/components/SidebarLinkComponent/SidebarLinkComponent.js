import {defineComponent} from 'vue'
import {useGlobalStore} from "stores/global";

export default defineComponent({
  name: 'EssentialLink',
  data() {
    return {
      globalStore: useGlobalStore(),
    }
  },
  props: {
    title: {
      type: String,
      required: true
    },

    caption: {
      type: String,
      default: ''
    },

    link: {
      type: String,
      default: '#'
    },

    icon: {
      type: String,
      default: ''
    },

    permission: {
      type: String,
      default: ''
    }
  }
})
