Feature: Singleton

  Scenario: A user add singleton object twice
    Given A user has FirstInterface interface
    And A user has FirstImpl struct that implement FirstInterface
    When A user add FirstImpl singleton to DI container twice
    Then No error occurs