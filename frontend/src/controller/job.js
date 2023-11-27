export default (client) => ({
  subscribeJob(subscriptionIdBefore, id, callback, error) {
    client.graphql.unsubscribe(subscriptionIdBefore);
    return client.graphql.subscribeTrackedObject(`query {
      data: getJob(id: "${id}") {
        id
        name
        strategy
        destinationStorageId
        originStorageId
        executorSatelliteId
        configuration
        createdAt
        updatedAt
        description
      }
    }`, (data, subscribeId) => callback(data, subscribeId), (data) => error(data.message));
  },
  subscribeJobs(subscriptionIdBefore, args, callback, error) {
    let search = args.search.replaceAll("\\", "\\\\").replaceAll('"', '\\"');

    client.graphql.unsubscribe(subscriptionIdBefore);
    return client.graphql.subscribeTrackedObject(`query {
      data: getJobs(limit: ${args.limit}, search: "${search}", page: ${args.page}) {
        rows  {
          id
          name
          strategy
          createdAt
          updatedAt
        }
        page {
          totalRows
          totalPages
        }
      }
    }`, (data, subscribeId) => callback(data, subscribeId), (data) => error(data.message));
  },
  subscribeJobStrategies(subscriptionIdBefore, callback, error) {
    client.graphql.unsubscribe(subscriptionIdBefore);
    return client.graphql.subscribeTrackedObject(`query {
      data: getJobStrategies {
        name
        configuration {
          name
          kind
          default
        }
      }
    }`, (data, subscribeId) => callback(data, subscribeId), (data) => error(data.message));
  },
  createJob(args, callback, error) {
    let fullQuery = client.graphql.buildQuery({
      method: "createJob",
      args,
      fields: [
        "id",
      ],
      isMutation: true
    });

    let updatedSubscriptionId = client.graphql.subscribeTrackedObject(
      fullQuery,
      (data) => {
        client.graphql.unsubscribe(updatedSubscriptionId);
        callback(data);
      },
      (data) => error(data.message)
    );
  },
  executeBackup(args, callback, error) {
    let fullQuery = client.graphql.buildQuery({
      method: "executeBackup",
      args,
      fields: [
        "id",
        "jobId",
      ],
      isMutation: true
    });

    let updatedSubscriptionId = client.graphql.subscribeTrackedObject(
      fullQuery,
      (data) => {
        client.graphql.unsubscribe(updatedSubscriptionId);
        callback(data);
      },
      (data) => error(data.message)
    );
  },
  updateJob(args, callback, error) {
    let fullQuery = client.graphql.buildQuery({
      method: "updateJob",
      args,
      fields: [
        "id",
      ],
      isMutation: true
    });

    let updatedSubscriptionId = client.graphql.subscribeTrackedObject(
      fullQuery,
      (data) => {
        client.graphql.unsubscribe(updatedSubscriptionId);
        callback(data);
      },
      (data) => error(data.message)
    );
  },
  deleteJob(args, callback, error) {
    let fullQuery = client.graphql.buildQuery({
      method: "deleteJob",
      args: {
        id: args.id,
      },
      fields: [
        "id",
      ],
      isMutation: true
    });

    let updatedSubscriptionId = client.graphql.subscribeTrackedObject(
      fullQuery,
      (data) => {
        client.graphql.unsubscribe(updatedSubscriptionId);
        callback(data);
      },
      (data) => error(data.message)
    );
  }
});
