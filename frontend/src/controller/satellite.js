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
          connected
        }
        page {
          totalRows
          totalPages
        }
      }
    }`, (data, subscribeId) => callback(data, subscribeId), (data) => error(data.message));
  },
  subscribeSatelliteLogs(subscriptionIdBefore, args, callback, error) {
    let search = args.search.replaceAll("\\", "\\\\").replaceAll('"', '\\"');

    client.graphql.unsubscribe(subscriptionIdBefore);
    return client.graphql.subscribeTrackedObject(`query {
      data: getSatelliteLogs(satelliteId: "${args.id}", limit: ${args.limit}, search: "${search}", page: ${args.page}) {
        rows {
          id
          sender
          message
          level
          loggedAt
        }
        page {
          totalRows
          totalPages
        }
      }
    }`, (data, subscribeId) => {
      data.rows = data.rows.sort().reverse();
      callback(data, subscribeId);
    }, (data) => error(data.message));
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
        connected
      }
    }`, (data, subscribeId) => callback(data, subscribeId), (data) => error(data.message));
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
  resyncSatellite(args, callback, error) {
    let fullQuery = client.graphql.buildQuery({
      method: "resyncSatellite",
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
