export default (client) => ({
  subscribeStorages(subscriptionIdBefore, args, callback, error) {
    let search = args.search.replaceAll("\\", "\\\\").replaceAll('"', '\\"');

    client.graphql.unsubscribe(subscriptionIdBefore);
    return client.graphql.subscribeTrackedObject(`query {
      data: getStorages(limit: ${args.limit}, search: "${search}", page: ${args.page}) {
        rows  {
          id
          name
          kind
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
  subscribeStorage(subscriptionIdBefore, id, callback, error) {
    client.graphql.unsubscribe(subscriptionIdBefore);
    return client.graphql.subscribeTrackedObject(`query {
      data: getStorage(id: "${id}") {
        id
        name
        kind
        configuration
        createdAt
        updatedAt
        description
      }
    }`, (data) => callback(data), (data) => error(data.message));
  },
  subscribeStorageKinds(subscriptionIdBefore, callback, error) {
    client.graphql.unsubscribe(subscriptionIdBefore);
    return client.graphql.subscribeTrackedObject(`query {
      data: getStorageKinds {
        name
        configuration {
          name
          kind
          default
        }
      }
    }`, (data) => callback(data), (data) => error(data.message));
  },
  createStorage(args, callback, error) {
    let fullQuery = client.graphql.buildQuery({
      method: "createStorage",
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
  updateStorage(args, callback, error) {
    let fullQuery = client.graphql.buildQuery({
      method: "updateStorage",
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
  deleteStorage(args, callback, error) {
    let fullQuery = client.graphql.buildQuery({
      method: "deleteStorage",
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
