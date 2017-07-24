angular.module('app')
    .controller('job', [
    '$scope',
    'job',
    function ($scope, jobService) {

        $scope.removeParam = function(p){
            for (var i=0;i<$scope.params.length;i++) {
                if ($scope.params[i] == p) {
                    $scope.params.splice(i, 1);
                    break;
                }
            }
        };

        $scope.removeBuildStep = function(s) {
            for (var i=0;i<$scope.buildSteps.length;i++) {
                if ($scope.buildSteps[i] == s) {
                    $scope.buildSteps.splice(i, 1);
                    break;
                }
            }
        };

        $scope.parametrized = false;
        $scope.params = [];
        $scope.buildSteps = [];

        $scope.addChoiceParam = function(){
            $scope.params.push({
                type: 'choice'
            });

            return false;
        };

        $scope.addStringParam = function(){
            $scope.params.push({
                type: 'string'
            });

            return false;
        };

        $scope.loadPluginSchema = function (pluginName, pluginDescription) {
            jobService.pluginSchema({name: pluginName}, function(data){
                $scope.buildSteps.push({
                    name: pluginName,
                    description: pluginDescription,
                    schema: data
                });
            });
        }


    }
])
    .factory('job', function($resource) {
        return $resource('/a/jobs', {}, {
            pluginSchema: {
                method: 'GET',
                url: '/a/jobs/plugins/schema',
                isArray: true
            }
        })
    });