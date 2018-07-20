import { Component, OnInit, ViewChild, Input, Inject, forwardRef  } from '@angular/core';
import { Comparator, State }                                        from "clarity-angular";

import { Area }                        from '../../../../model/time/area';
import { Resource }                    from '../../../../model/time/resource';
import { Phase }                       from '../../../../model/time/phase';
import { SessionUser }                 from '../../../../model/base/sign';
import { UserResource }                from '../../../../model/time/user-resource';
import { AreaService }                 from '../../../../service/time/area.service';
import { PhaseService }                from '../../../../service/time/phase.service';
import { UserResourceService }         from '../../../../service/time/user-resource.service';
import { SignService }                 from '../../../../service/system/sign.service';
import { CustomComparator, doSorting } from '../../../../shared/utils';

import { TaskListComponent } from '../../task/list/list.component';
import { ShellComponent }    from '../../../../base/shell/shell.component';

@Component({
  selector: 'time-user-resource-list',
  templateUrl: './list.component.html',
  styleUrls: ['./list.component.css']
})
export class UserResourceListComponent implements OnInit {
  @ViewChild(TaskListComponent)
  taskListComponent: TaskListComponent;

  userResources: UserResource[];
  fieldPhasesMap: Map<number, Phase[]> = new Map();

  @Input() initAllFlag: boolean = false;
  filterArea: Area = new Area();
  filterAreaSumTime: number = 0;
  filterAreaPhase: Phase = new Phase();
  currentState: State;

  timeComparator: Comparator<UserResource> = new CustomComparator<UserResource>("Time", "number");

  pageSize: number = 10;
  totalCount: number = 0;
  currentPage: number = 1;

  constructor(
    private areaService: AreaService,
    private phaseService: PhaseService,
    private userResourceService: UserResourceService,
    private signService: SignService,
    @Inject(forwardRef(() => ShellComponent))
    private shell: ShellComponent,
  ) { }

  ngOnInit() {
    this.pageSize = 10;
    this.phaseService.List(null).subscribe(res => {
      res.forEach((one, k) => {
        if (this.fieldPhasesMap.get(one.FieldId) == undefined) {
          this.fieldPhasesMap.set(one.FieldId, []);
        }
        let phaseArray = this.fieldPhasesMap.get(one.FieldId);
        phaseArray.push(one);
        this.fieldPhasesMap.set(one.FieldId, phaseArray);
      })
    })
  }

  changeFilterArea(areaId: number) {
    this.areaService.Get(areaId).subscribe(area => {
      this.filterArea = area;
      this.refreshAll(0, this.pageSize);
    })
  }

  load(state: any): void {
    this.currentState = state;
    if (this.currentState && this.currentState.sort) {
      this.refreshAll(0, this.pageSize)
    }
  }

  refresh() {
    this.currentPage = 1;
    this.refreshAll(0, this.pageSize);
  }

  refreshAll(from: number, to: number): void {
    if ((this.filterArea.Id == undefined) && (this.initAllFlag)) {
      this.userResourceService.List(null).subscribe(res => {
        this.totalCount = res.length;
        this.userResources = res.slice(from, to);
      })
    } else if (this.filterArea.Id != undefined) {
      let userResource = new UserResource();
      userResource.UserId = this.shell.currentUser.Id;
      userResource.Resource.Area = new Area();
      userResource.Resource.Area.Id = this.filterArea.Id;
      userResource.Resource.WithSub = true;
      this.userResourceService.List(userResource).subscribe(res => {
        this.filterAreaSumTime = 0;
        res.forEach((one, k) => {
          this.filterAreaSumTime += one.Time;
        })
        this.filterAreaSumTime = this.filterAreaSumTime / 60;
        let relateFieldPhaseArray = this.fieldPhasesMap.get(this.filterArea.FieldId);
        let tmpPhase = new Phase();
        for (let k in relateFieldPhaseArray) {
          if (relateFieldPhaseArray[k].Threshold < this.filterAreaSumTime) {
            continue
          }
          if (tmpPhase.Id == undefined) {
            tmpPhase = relateFieldPhaseArray[k];
            continue;
          }
          if (relateFieldPhaseArray[k].Threshold <= tmpPhase.Threshold) {
            tmpPhase = relateFieldPhaseArray[k];
          }
        }
        this.filterAreaPhase = tmpPhase;
        this.totalCount = res.length;
        if (this.currentState && this.currentState.sort) {
          res = doSorting<UserResource>(res, this.currentState);
        }
        this.userResources = res.slice(from, to);
      })
    }
  }

  listTasks(userResource: UserResource): void {
    this.taskListComponent.changeFilterResource(userResource.Resource);
  }
}
