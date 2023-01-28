import utils from './utils';

describe('Validate Permissions for Help Desk User', () => {

    beforeEach(() => {
        cy.fixture('users').then((users) => {
            cy.fixture('base').then((base) => {
                cy.login(base.url, users.users[3].name, users.users[3].password)
            });
        });
    })

    utils.validateUserListViewable()
    utils.validateUserCreationButtonNotToBeVisible()
    utils.validateUserCreationIsBlockedByBackend()
    utils.validatePermissionUpdateIsNotAllowed()
    utils.validateUserUpdateIsNotAllowed()

})