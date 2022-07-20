using Dapper;
using System.Data;
using System.Data.SqlClient;
using Microsoft.Extensions.Configuration;
using Npgsql;

public class Database
{
  private readonly IConfiguration _config;

  public Database(IConfiguration config)
  {
    _config = config;
  }

  private string connectionString
  {
    get => _config.GetConnectionString("Default");
  }

  public async Task<IEnumerable<T>> LoadData<T>(string query)
  {
    using IDbConnection connection = new NpgsqlConnection(connectionString);

    return await connection.QueryAsync<T>(query);
  }

  public async Task SaveData(string query)
  {
    using IDbConnection connection = new NpgsqlConnection(connectionString);

    await connection.ExecuteAsync(query);
  }
}
