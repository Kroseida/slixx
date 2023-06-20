import {boot} from 'quasar/wrappers'
import controller from "../controller";

export default boot(({app, store}) => {
  const controllerObj = controller()
  controllerObj.connect()

  app.config.globalProperties.$controller = controllerObj

  store.use(() => ({$controller: controllerObj}));
})

export {}
