angular.module('app')
    .controller('job', [
    '$scope',
    'job',
    'cfpLoadingBar',
    '$timeout',
    function ($scope, jobService,cfpLoadingBar,$timeout) {
        $scope.applied = false;

        $scope.job = {
            is_enabled: true,
            params : [],
            steps: []
        };

        if(document.location.pathname.match(/edit\/\d+/)) {
            $scope.job["id"] = parseInt(document.location.pathname.match(/edit\/(\d+)/)[1]);
            jobService.get({id: $scope.job["id"]}, function(data){
                //console.log(data);
                $scope.job = data;
                if ($scope.job.params && $scope.job.params.length > 0) {
                    $scope.parametrized = true;
                }
                if ($scope.job.remote_token && $scope.job.remote_token != "") {
                    $scope.buildRemotely = true;
                }
                if ($scope.job.build_path && $scope.job.build_path != "") {
                    $scope.buildPath = true;
                }
                if ($scope.job.schedule && $scope.job.schedule != "") {
                    $scope.buildPeriodically = true;
                }
                if(!$scope.job.params) {
                    $scope.job.params = [];
                }
                if(!$scope.job.steps) {
                    $scope.job.steps = [];
                }
            });
        }

        $scope.removeParam = function(p){
            for (var i=0;i<$scope.job.params.length;i++) {
                if ($scope.job.params[i] == p) {
                    $scope.job.params.splice(i, 1);
                    break;
                }
            }
        };



        $scope.removeBuildStep = function(s) {
            for (var i=0;i<$scope.job.steps.length;i++) {
                if ($scope.job.steps[i] == s) {
                    $scope.job.steps.splice(i, 1);
                    break;
                }
            }
        };

        // $scope.params = [];
        // $scope.buildSteps = [];

        $scope.addChoiceParam = function(){
            $scope.job.params.push({
                type: 'choice'
            });

            return false;
        };

        $scope.addStringParam = function(){
            $scope.job.params.push({
                type: 'string'
            });

            return false;
        };

        $scope.loadPluginSchema = function (pluginName, pluginDescription) {
            jobService.pluginSchema({name: pluginName}, function(data){
                $scope.job.steps.push({
                    name: pluginName,
                    description: pluginDescription,
                    schema: data
                });
            });
        };
        $scope.formSubmit = function(){
            jobService.create($scope.job, function(){
                document.location = "/a/jobs";
            });
        };

        $scope.apply = function(){
            cfpLoadingBar.start();
            jobService.create($scope.job, function(data){
                $scope.job = data;
                cfpLoadingBar.complete();
                $scope.applied = true;

                $timeout(function(){
                    $scope.applied = false;
                }, 2000);
                // setTimeout(function(){}, 1000);
            });
        };

    }
])
    .controller('build', ['$scope', function($scope){
        console.log($scope);
    }])
    .factory('job', function($resource) {
        return $resource('/a/jobs', {}, {
            pluginSchema: {
                method: 'GET',
                url: '/a/jobs/plugins/schema/:name',
                isArray: true
            },
            create: {
                method: 'POST',
                url: '/a/jobs/create'
            },
            get: {
                method: 'GET',
                url: '/a/jobs/edit/:id',
                headers: { 'Accept': 'application/json' }
            }
        })
    });