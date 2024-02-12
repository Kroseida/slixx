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
});
