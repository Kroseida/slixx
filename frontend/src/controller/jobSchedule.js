export default (client) => ({
  subscribeJobSchedules(subscriptionIdBefore, args, callback, error) {
    let search = args.search.replaceAll("\\", "\\\\").replaceAll('"', '\\"');
    let sort = args.sort || ""

    client.graphql.unsubscribe(subscriptionIdBefore);
    if (args.jobId) {
      return client.graphql.subscribeTrackedObject(`query {
        data: getJobSchedules(limit: ${args.limit}, search: "${search}", page: ${args.page}, jobId: "${args.jobId}", sort: "${sort}") {
          rows  {
            id
            name
            description
            kind
            configuration
            createdAt
            updatedAt
          }
          page {
            totalRows
            totalPages
          }
        }
      }`, (data, subscribeId) => callback(data, subscribeId), (data) => error(data.message));
    }

    return client.graphql.subscribeTrackedObject(`query {
      data: getJobSchedules(limit: ${args.limit}, search: "${search}", page: ${args.page}, sort: "${sort}") {
        rows  {
            id
            name
            description
            kind
            configuration
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
  subscribeJobScheduleKinds(subscriptionIdBefore, callback, error) {
    client.graphql.unsubscribe(subscriptionIdBefore);
    return client.graphql.subscribeTrackedObject(`query {
      data: getJobScheduleKinds {
        name
        configuration {
          name
          kind
          default
        }
      }
    }`, (data, subscribeId) => callback(data, subscribeId), (data) => error(data.message));
  },
  createJobSchedule(args, callback, error) {
    let fullQuery = client.graphql.buildQuery({
      method: "createJobSchedule",
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
  updateJobSchedule(args, callback, error) {
    let fullQuery = client.graphql.buildQuery({
      method: "updateJobSchedule",
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
  deleteJobSchedule(args, callback, error) {
    let fullQuery = client.graphql.buildQuery({
      method: "deleteJobSchedule",
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
