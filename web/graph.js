fetch('/data.json')
  .then(response => response.json())
  .then(data => chart(data));

let highPoints = 70
let middlePoints = 50
let lowPoints = 20
function chart(json) {
    let labels = [];
    let data = [];
    let colors = [];

    json.forEach(function (item, index) {
        labels.push(item.Family)
        data.push(item.AveragePoints)
        if (item.AveragePoints > highPoints) {
            colors.push('rgba(255, 99, 132, 0.5)')
        } else if (item.AveragePoints > middlePoints) {
            colors.push('rgba(255, 159, 64, 0.4)')
        } else if (item.AveragePoints > lowPoints) {
            colors.push('rgba(75, 192, 192, 0.5)')
        } else {
            colors.push('rgba(0, 0, 0, 0.2)')
        }
    });
    console.log(labels)
    var ctx = document.getElementById('myChart').getContext('2d');
    var myChart = new Chart(ctx, {
        type: 'bar',
        data: {
            labels,
            datasets: [{
                label: 'Average Points',
                data,
                backgroundColor: colors,
                borderWidth: 2
            }]
        },
        options: {
            scales: {
                yAxes: [{
                    ticks: {
                        beginAtZero: true
                    }
                }]
            },
            responsive: true,
            maintainAspectRatio: false,
        }
    });
}
