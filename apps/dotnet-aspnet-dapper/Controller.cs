using Microsoft.AspNetCore.Mvc;

[ApiController]
[Route("/")]
public class WeatherForecastController : ControllerBase
{

  private readonly Database _database;

  public WeatherForecastController(Database database)
  {
    _database = database;
  }

  [HttpGet]
  public async Task<User> Get()
  {
    // Create
    var createQuery = @"
    INSERT INTO users (name, age)
      VALUES('Test!', 1) RETURNING id;
    ";

    var userId = (await _database.LoadData<int>(createQuery)).First();

    // Read
    var readQuery = @$"
      SELECT id, name, age
      FROM USERS
      WHERE id = {userId}
    ";

    var user = (await _database.LoadData<User>(readQuery)).First();

    // Update
    user.Age += 1;

    var updateQuery = @$"
      UPDATE users SET
        name = '{user.Name}',
        age  = '{user.Age}'
      WHERE id = {user.Id};
    ";

    await _database.SaveData(updateQuery);

    // Delete
    var deleteQuery = @$"DELETE FROM users WHERE id = {user.Id}";
    await _database.SaveData(updateQuery);

    return user;
  }
}
