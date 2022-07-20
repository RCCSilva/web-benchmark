using Microsoft.AspNetCore.Mvc;

[ApiController]
[Route("/")]
public class WeatherForecastController : ControllerBase
{

  private readonly UserContext _userContext;

  public WeatherForecastController(UserContext userContext)
  {
    _userContext = userContext;
  }

  [HttpGet]
  public async Task<User> Get()
  {
    // Create
    try
    {

      var user = new User
      {
        Name = "Test!",
        Age = 1
      };
      await _userContext.AddAsync(user);
      await _userContext.SaveChangesAsync();

      // Read
      var userDb = _userContext.Users.Where(u => u.Id == user.Id).SingleOrDefault()!;

      // Update
      userDb.Age += 1;
      await _userContext.SaveChangesAsync();

      // Delete
      _userContext.Remove(userDb);

      return userDb;
    } catch (Exception e) {
      Console.WriteLine("Failed {0}", e);
      throw e;
    }
  }
}
