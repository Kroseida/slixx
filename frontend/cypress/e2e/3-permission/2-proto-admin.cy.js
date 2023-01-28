import utils from './utils';

describe('Validate Permissions for Prototype Admin', () => {

    beforeEach(() => {
        cy.fixture('users').then((users) => {
            cy.fixture('base').then((base) => {
                cy.login(base.url, users.users[1].name, users.users[1].password)
            });
        });
    })

    utils.validateUserListViewable()
    utils.validateUserCreationButtonNotToBeVisible()
    utils.validateUserUpdateIsAllowed()
    utils.validateUserCreationIsBlockedByBackend()

    it('Validate permission update is allowed', () => {
        cy.get('#sidebar__user').click()

        cy.get('table')
            .find('tbody')
            .find('tr')
            .first()
            .find('td')
            .first()
            .click()

        cy.get('#permission__user_view').click()
        cy.get('.success').should('be.visible')
    })
})