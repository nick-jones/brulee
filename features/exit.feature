Feature: Exit statement

  Scenario: Exit alone
    Given the program:
    """
    exit
    score(x) = 1
    """
    When the program is run
    Then the score output is empty

  Scenario: Exit within conditional
    Given the program:
    """
    when
      "x" == "x"
    then
      exit
    done
    score(x) = 1
    """
    When the program is run
    Then the score output is empty

  Scenario: Exit within conditional where the condition fails
    Given the program:
    """
    when
      "x" == "y"
    then
      exit
    done
    score(x) = 1
    """
    When the program is run
    Then the score output is:
      | Name | Score |
      | x    | 1     |
