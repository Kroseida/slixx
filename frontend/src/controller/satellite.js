export default (client) => ({
  subscribeSatellites(subscriptionIdBefore, args, callback, error) {
    let search = args.search.replaceAll("\\", "\\\\").replaceAll('"', '\\"');

    client.graphql.unsubscribe(subscriptionIdBefore);
    return client.graphql.subscribeTrackedObject(`query {
      data: getSatellites(limit: ${args.limit}, search: "${search}", page: ${args.page}) {
        rows  {
          id
          name
          address
          createdAt
          updatedAt
        }
        page {
          totalRows
          totalPages
        }
      }
    }`, (data) => callback(data), (data) => error(data.message));
  },
  subscribeSatellite(subscriptionIdBefore, id, callback, error) {
    client.graphql.unsubscribe(subscriptionIdBefore);
    return client.graphql.subscribeTrackedObject(`query {
      data: getSatellite(id: "${id}") {
        id
        name
        address
        createdAt
        updatedAt
        description
      }
    }`, (data) => callback(data), (data) => error(data.message));
  },
  deleteSatellite(args, callback, error) {
    let fullQuery = client.graphql.buildQuery({
      method: "deleteSatellite",
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
  },
  createSatellite(args, callback, error) {
    let fullQuery = client.graphql.buildQuery({
      method: "createSatellite",
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
  updateSatellite(args, callback, error) {
    let fullQuery = client.graphql.buildQuery({
      method: "updateSatellite",
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
