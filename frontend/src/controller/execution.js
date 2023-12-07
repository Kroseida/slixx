export default (client) => ({
  subscribeExecutions(subscriptionIdBefore, args, callback, error) {
    let search = args.search.replaceAll("\\", "\\\\").replaceAll('"', '\\"');

    client.graphql.unsubscribe(subscriptionIdBefore);
    if (args.jobId) {
      return client.graphql.subscribeTrackedObject(`query {
        data: getExecutions(limit: ${args.limit}, search: "${search}", page: ${args.page}, jobId: "${args.jobId}") {
          rows  {
            id
            jobId
            kind
            createdAt
            finishedAt
            updatedAt
            status
          }
          page {
            totalRows
            totalPages
          }
        }
      }`, (data, subscribeId) => callback(data, subscribeId), (data) => error(data.message));
    }

    return client.graphql.subscribeTrackedObject(`query {
      data: getExecutions(limit: ${args.limit}, search: "${search}", page: ${args.page}) {
        rows  {
          id
          jobId
          createdAt
          updatedAt
          status
        }
        page {
          totalRows
          totalPages
        }
      }
    }`, (data, subscribeId) => callback(data, subscribeId), (data) => error(data.message));
  },
  subscribeExecutionHistory(subscriptionIdBefore, args, callback, error) {
    client.graphql.unsubscribe(subscriptionIdBefore);
    return client.graphql.subscribeTrackedObject(`query {
      data: getExecutionHistory(executionId: "${args.executionId}") {
        id
        message
        statusType
        percentage
        createdAt
      }
    }`, (data, subscribeId) => callback(data, subscribeId), (data) => error(data.message));
  }
});
