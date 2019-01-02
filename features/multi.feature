Feature:

  Scenario: 2 levels
    Given the program:
    """
    when
      "x" == "x"
    then
      score(x) = 1
    done

    when
      "y" == "y"
    then
      score(x) = 2
    done
    """
    When the program is run
    Then the score output is:
      | Name | Score |
      | x    | 2     |