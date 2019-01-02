Feature:

  Scenario: 2 levels
    Given the program:
    """
    when
      "x" == "x"
    then
      when
        "y" == "y"
      then
         score(x) = 1
      done
    done
    """
    When the program is run
    Then the score output is:
      | Name | Score |
      | x    | 1     |

  Scenario: 2 levels, first condition failing
    Given the program:
    """
    when
      "x" == "y"
    then
      when
        "y" == "y"
      then
         score(x) = 1
      done
    done
    """
    When the program is run
    Then the score output is empty

  Scenario: 2 levels, second condition failing
    Given the program:
    """
    when
      "x" == "x"
    then
      when
        "y" == "z"
      then
         score(x) = 1
      done
    done
    """
    When the program is run
    Then the score output is empty

  Scenario: 3 levels
    Given the program:
    """
    when
      "x" == "x"
    then
      when
        "y" == "y"
      then
        when
          "z" == "z"
        then
          score(x) = 1
        done
      done
    done
    """
    When the program is run
    Then the score output is:
      | Name | Score |
      | x    | 1     |
