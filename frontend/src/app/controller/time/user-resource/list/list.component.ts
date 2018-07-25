import { Component, OnInit, ViewChild, Input } from '@angular/core';
import { Inject, forwardRef }                  from '@angular/core';
import { Comparator, State, SortOrder}         from "clarity-angular";

import { PageSet }             from '../../../../model/base/basic';
import { Area }                from '../../../../model/time/area';
import { Resource }            from '../../../../model/time/resource';
import { Phase }               from '../../../../model/time/phase';
import { SessionUser }         from '../../../../model/base/sign';
import { UserResource }        from '../../../../model/time/user-resource';
import { AreaService }         from '../../../../service/time/area.service';
import { PhaseService }        from '../../../../service/time/phase.service';
import { UserResourceService } from '../../../../service/time/user-resource.service';
import { SignService }         from '../../../../service/system/sign.service';

import { CustomComparator }                             from '../../../../shared/utils';
import { loadPageFilterSort, reloadState, deleteState } from '../../../../shared/utils';
import { PageSize }                                     from '../../../../shared/const';

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

  timeComparator: Comparator<UserResource> = new CustomComparator<UserResource>("Time", "number");
  preSorted = SortOrder.Desc;

  filterArea: Area = new Area();
  filterAreaLeafIds: number [] = [];
  filterAreaSumTime: number = 0;
  filterAreaPhase: Phase = new Phase();

  currentState: State;
  pageSet: PageSet = new PageSet();

  constructor(
    private areaService: AreaService,
    private phaseService: PhaseService,
    private userResourceService: UserResourceService,
    private signService: SignService,
    @Inject(forwardRef(() => ShellComponent))
    private shell: ShellComponent,
  ) { }

  ngOnInit() {
    this.pageSet.Size = PageSize.Normal;
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

  refresh() {
    let state = reloadState(this.currentState, this.pageSet);
    this.load(state);
  }

  changeFilterArea(areaId: number) {
    this.areaService.Get(areaId).subscribe(area => {
      this.filterArea = area;
      this.areaService.ListAllLeaf(this.filterArea.Id).subscribe(areaIds => {
        areaIds.push(this.filterArea.Id);
        this.filterAreaLeafIds = areaIds;
        this.refresh();
      })
    })
  }

  load(state: State): void {
    if ((this.filterArea.Id == undefined) && (this.initAllFlag)) {
      let userResource = new UserResource();
      userResource = loadPageFilterSort<UserResource>(userResource, state);
      this.pageSet.Current = userResource.Page.Current;
      this.currentState = state;
      this.userResourceService.Count(userResource).subscribe(count => {
        this.pageSet.Count = count;
        this.userResourceService.List(userResource).subscribe(res => {
          this.userResources = res;
        })
      })
    } else if (this.filterArea.Id != undefined) {
      let userResource = new UserResource();
      userResource.UserId = this.shell.currentUser.Id;
      userResource.Resource.Area = new Area();
      userResource.Resource.Area.Ids = this.filterAreaLeafIds;
      userResource = loadPageFilterSort<UserResource>(userResource, state);
      this.pageSet.Current = userResource.Page.Current;
      this.currentState = state;
      this.userResourceService.Count(userResource).subscribe(count => {
        this.pageSet.Count = count;
        this.userResourceService.List(userResource).subscribe(res => {
          this.userResourceService.AreaSumTime(userResource).subscribe(sumTime => {
            this.filterAreaSumTime = sumTime / 60;
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
            this.userResources = res;
          })
        })
      })
    }
  }

  listTasks(userResource: UserResource): void {
    this.taskListComponent.changeFilterResource(userResource.Resource);
  }

}
