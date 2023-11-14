import graphql from "./graphql";
import user from "./user";
import job from "./job";
import storage from "./storage";
import satellite from "./satellite";
import backup from "./backup";

export default () => ({
  graphql: null,
  user: null,
  job: null,
  storage: null,
  satellite: null,
  backup: null,
  connect() {
    this.graphql = graphql.createClient(localStorage.getItem('_auth'));
    this.user = user(this);
    this.job = job(this);
    this.storage = storage(this);
    this.satellite = satellite(this);
    this.backup = backup(this);
  }
})
