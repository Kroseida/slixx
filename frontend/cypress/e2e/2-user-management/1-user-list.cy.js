describe('View the user list', () => {

    beforeEach(() => {
        cy.loginAdmin()
    })

    it('check validity of user list', () => {
        cy.get('#sidebar__user').click()

        cy.get('table').should('be.visible')

        cy.get('table')
            .find('tbody')
            .find('tr')
            .should('have.length', 1)

        cy.get('table')
            .find('tbody')
            .find('tr')
            .find('td')
            .eq(0)
            .should('contain', 'admin')
    })

});