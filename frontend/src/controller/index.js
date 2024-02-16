import graphql from "./graphql";
import user from "./user";
import job from "./job";
import storage from "./storage";
import satellite from "./satellite";
import backup from "./backup";
import execution from "./execution";
import jobSchedule from "./jobSchedule";

export default () => ({
  graphql: null,
  user: null,
  job: null,
  storage: null,
  satellite: null,
  backup: null,
  execution: null,
  jobSchedule: null,
  connect() {
    this.graphql = graphql.createClient(localStorage.getItem('_auth'));
    this.user = user(this);
    this.job = job(this);
    this.storage = storage(this);
    this.satellite = satellite(this);
    this.backup = backup(this);
    this.execution = execution(this);
    this.jobSchedule = jobSchedule(this);
  },
  unsubscribe(id) {
    this.graphql.unsubscribe(id);
  }
})
