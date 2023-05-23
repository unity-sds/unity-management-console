<script>
import { goto } from "$app/navigation";
</script>
<div class="container">
    <h2>SPS Installation Information</h2>
<div class="accordion" id="accordionExample">
  <div class="accordion-item">
    <h2 class="accordion-header" id="headingOne">
      <button class="accordion-button" type="button" data-bs-toggle="collapse" data-bs-target="#collapseOne" aria-expanded="true" aria-controls="collapseOne">
        Info
      </button>
    </h2>
    <div id="collapseOne" class="accordion-collapse collapse show" aria-labelledby="headingOne" data-bs-parent="#accordionExample">
        <div class="accordion-body">
            <table class="table">
                <tbody>
                    <tr>
                        <th scope="row">Application Name</th>
                        <td>SPS</td>
                    </tr>
                    <tr>
                        <th scope="row">Application Version</th>
                        <td>1.0.0</td>
                    </tr>
                    <tr>
                        <th scope="row">Application Install Date</th>
                        <td>01/01/2020</td>
                    </tr>
                    <tr>
                        <th scope="row">Application Owner</th>
                        <td>Tom Barber</td>
                    </tr>
                </tbody>
            </table>
        </div>
    </div>
  </div>
  <div class="accordion-item">
    <h2 class="accordion-header" id="headingTwo">
      <button class="accordion-button collapsed" type="button" data-bs-toggle="collapse" data-bs-target="#collapseTwo" aria-expanded="false" aria-controls="collapseTwo">
        Access
      </button>
    </h2>
    <div id="collapseTwo" class="accordion-collapse collapse" aria-labelledby="headingTwo" data-bs-parent="#accordionExample">
        <div class="accordion-body">
            <table class="table">
                <tbody>
                    <tr>
                        <th scope="row">Mozart URL</th>
                        <td><a href="http://unity.jpl.nasa.gov/myproject/sps/mozart/">http://unity.jpl.nasa.gov/myproject/sps/mozart/</a></td>
                    </tr>
                </tbody>
            </table>
        </div>
    </div>
  </div>
  <div class="accordion-item">
      <h2 class="accordion-header" id="headingFive">
          <button class="accordion-button collapsed" type="button" data-bs-toggle="collapse" data-bs-target="#collapseFive" aria-expanded="false" aria-controls="collapseFive">
              Extensions
          </button>
      </h2>
      <div id="collapseFive" class="accordion-collapse collapse" aria-labelledby="headingFive" data-bs-parent="#accordionExample">
          <div class="accordion-body">
              <div class="card">
                  <div class="card-body">
                      <h5 class="card-title">Kubernetes</h5>
                      <p class="card-text">SPS uses the Kubernetes extension.</p>
                  </div>
                  <div class="card-footer">
                      <button class="btn btn-primary" on:click={() => goto('/ui/extension/kubernetes/cluster')}>Manage</button>
                  </div>
              </div>
          </div>
      </div>
  </div>
  <div class="accordion-item">
    <h2 class="accordion-header" id="headingThree">
      <button class="accordion-button collapsed" type="button" data-bs-toggle="collapse" data-bs-target="#collapseThree" aria-expanded="false" aria-controls="collapseThree">
        Tags
      </button>
    </h2>
<div id="collapseThree" class="accordion-collapse collapse" aria-labelledby="headingThree" data-bs-parent="#accordionExample">
  <div class="accordion-body">
    <table class="table">
      <tbody>
        <tr>
          <th scope="row">Venue</th>
          <td>My Venue</td>
        </tr>
        <tr>
          <th scope="row">ServiceArea</th>
          <td>SPS</td>
        </tr>
        <tr>
          <th scope="row">CapVersion</th>
          <td>1.0.0</td>
        </tr>
        <tr>
          <th scope="row">Component</th>
          <td>HySDS</td>
        </tr>
        <tr>
          <th scope="row">Name</th>
          <td>SPS</td>
        </tr>
        <tr>
          <th scope="row">Proj</th>
          <td>sips</td>
        </tr>
        <tr>
          <th scope="row">CreatedBy</th>
          <td>SPS</td>
        </tr>
        <tr>
          <th scope="row">Env</th>
          <td>My Venue</td>
        </tr>
        <tr>
          <th scope="row">Stack</th>
          <td>HySDS</td>
        </tr>
        <tr>
          <th scope="row">Owner</th>
          <td>Tom Barber</td>
        </tr>
        <tr>
          <th scope="row">ExposedWeb</th>
          <td>True</td>
        </tr>
        <tr>
          <th scope="row">UserFacing</th>
          <td>true</td>
        </tr>
        <tr>
          <th scope="row">POC</th>
          <td><a href="mailto:tom.barber@jpl.nasa.gov">tom.barber@jpl.nasa.gov</a></td>
        </tr>
      </tbody>
    </table>
  </div>
</div>
  </div>
  <div class="accordion-item">
    <h2 class="accordion-header" id="headingFour">
      <button class="accordion-button collapsed" type="button" data-bs-toggle="collapse" data-bs-target="#collapseFour" aria-expanded="false" aria-controls="collapseFour">
        Operation
      </button>
    </h2>
    <div id="collapseFour" class="accordion-collapse collapse" aria-labelledby="headingFour" data-bs-parent="#accordionExample">
        <div class="accordion-body">
            <div class="row">
                <div class="col-sm-6">
                    <div class="card">
                        <div class="card-body">
                            <h5 class="card-title">Upgrade</h5>
                            <p class="card-text">Upgrade the installed SPS platform to the latest version.</p>
                        </div>
                        <div class="card-footer">
                            <button class="btn btn-primary">Go</button>
                        </div>
                    </div>
                </div>
                <div class="col-sm-6">
                    <div class="card">
                        <div class="card-body">
                            <h5 class="card-title">Uninstall</h5>
                            <p class="card-text">Uninstall SPS completely from this project.</p>
                        </div>
                        <div class="card-footer">
                            <button class="btn btn-danger">Go</button>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
  </div>
</div>
<button class="btn btn-primary">Back</button>
</div>
