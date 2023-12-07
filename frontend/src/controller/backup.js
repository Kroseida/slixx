export default (client) => ({
  subscribeBackups(subscriptionIdBefore, args, callback, error) {
    let search = args.search.replaceAll("\\", "\\\\").replaceAll('"', '\\"');

    client.graphql.unsubscribe(subscriptionIdBefore);
    if (args.jobId) {
      return client.graphql.subscribeTrackedObject(`query {
        data: getBackups(limit: ${args.limit}, search: "${search}", page: ${args.page}, jobId: "${args.jobId}") {
          rows  {
            id
            name
            createdAt
            updatedAt
            jobId
          }
          page {
            totalRows
            totalPages
          }
        }
      }`, (data, subscribeId) => callback(data, subscribeId), (data) => error(data.message));
    }

    return client.graphql.subscribeTrackedObject(`query {
      data: getBackups(limit: ${args.limit}, search: "${search}", page: ${args.page}) {
        rows  {
          id
          name
          createdAt
          updatedAt
          jobId
        }
        page {
          totalRows
          totalPages
        }
      }
    }`, (data, subscribeId) => callback(data, subscribeId), (data) => error(data.message));
  },
  subscribeBackup(subscriptionIdBefore, id, callback, error) {
    client.graphql.unsubscribe(subscriptionIdBefore);
    return client.graphql.subscribeTrackedObject(`query {
      data: getBackup(id: "${id}") {
        id
        executionId
        createdAt
        updatedAt
      }
    }`, (data, subscribeId) => callback(data, subscribeId), (data) => error(data.message));
  },
  restoreBackup(args, callback, error) {
    let fullQuery = client.graphql.buildQuery({
      method: "restoreBackup",
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
});
