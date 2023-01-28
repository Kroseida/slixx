export default {
    validateUserListViewable() {
        it('Validate user list is viewable', () => {
            cy.get('#sidebar__user').should('be.visible')
        });
    },
    validateUserListNotViewable() {
        it('Validate user list to be not viewable', () => {
            cy.get('#sidebar__user').should('not.exist')
        });
    },
    validateUserCreationButtonNotToBeVisible() {
        it('Validate user creation button not to be visible', () => {
            cy.get('#sidebar__user').click()
            cy.get('#action__create__user').should('not.exist')
        });
    },
    validateUserUpdateIsAllowed() {
        it('Validate user update is allowed', () => {
            cy.get('#sidebar__user').click()

            cy.get('table')
                .find('tbody')
                .find('tr')
                .first()
                .find('td')
                .first()
                .click()

            cy.get('#userDetails__description').type('this is an update by prototype admin')
            cy.get('#userDetails__save__button').click()
            cy.get('.success').should('be.visible')
        });
    },
    validateUserUpdateIsNotAllowed() {
        it('Validate user update is not allowed', () => {
            cy.get('#sidebar__user').click()

            cy.get('table')
                .find('tbody')
                .find('tr')
                .first()
                .find('td')
                .first()
                .click()

            cy.get('#userDetails__name').should('have.attr', 'readonly', 'readonly')
            cy.get('#userDetails__firstName').should('have.attr', 'readonly', 'readonly')
            cy.get('#userDetails__lastName').should('have.attr', 'readonly', 'readonly')
            cy.get('#userDetails__email').should('have.attr', 'readonly', 'readonly')
            cy.get('#userDetails__description').should('have.attr', 'readonly', 'readonly')
        })
    },
    validateUserCreationIsBlockedByBackend() {
        it('Validate user creation is blocked by backend', () => {
            cy.fixture('base').then((base) => {
                cy.visit(base.url + '/#/app/user/new/')

                cy.get('#userDetails__name').type('test_user_of_prototype_admin')
                cy.get('#userDetails__firstName').type('firstName')
                cy.get('#userDetails__lastName').type('lastName')
                cy.get('#userDetails__email').type('mert@mert.de')
                cy.get('#userDetails__active').select('Active')
                cy.get('#userDetails__description').type('Max')

                cy.get('#userDetails__save__button').click()
                cy.get('.error').should('be.visible')
            });
        });

    },
    validatePermissionUpdateIsNotAllowed() {
        it('Validate permission update is not allowed', () => {
            cy.get('#sidebar__user').click()

            cy.get('table')
                .find('tbody')
                .find('tr')
                .first()
                .find('td')
                .first()
                .click()

            cy.get('#permission__user_view').should('not.exist')
        })
    }
}