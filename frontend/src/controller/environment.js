export default (client) => ({
  environment(subscriptionIdBefore, callback, error) {
    client.graphql.unsubscribe(subscriptionIdBefore);
    return client.graphql.subscribeTrackedObject(`query {
      data: environment {
        version
      }
    }`, (data, subscribeId) => callback(data, subscribeId), (data) => error(data.message));
  },
})
