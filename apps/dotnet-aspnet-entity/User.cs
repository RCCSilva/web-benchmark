using System.Collections.Generic;
using Microsoft.EntityFrameworkCore;
using System.ComponentModel.DataAnnotations.Schema;

public class UserContext : DbContext
{
  private readonly IConfiguration _configuration;

  public UserContext(IConfiguration configuration)
  {
    _configuration = configuration;
  }

  public DbSet<User> Users { get; set; }

  protected override void OnConfiguring(DbContextOptionsBuilder optionsBuilder)
      => optionsBuilder.UseNpgsql(_configuration.GetConnectionString("default"));
}

[Table("users")]
public class User
{

  [Column("id")]
  public int Id { get; set; }

  [Column("name")]
  public string? Name { get; set; }

  [Column("age")]
  public int Age { get; set; }
}
